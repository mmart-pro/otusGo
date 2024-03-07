package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	errs "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/errors"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/mqclient"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/service/events"
	memorystorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/pkg/logger"
)

type Scheduler struct {
	cfg config.SchedulerConfig
}

func NewScheduler(cfg config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		cfg: cfg,
	}
}

func (app Scheduler) Startup() error {
	// logger
	logg, err := logger.NewLogger(app.cfg.LoggerConfig.Level, app.cfg.LoggerConfig.LogFile)
	if err != nil {
		return fmt.Errorf("can't start logg: %w", err)
	}
	defer logg.Close()

	// IRL exclude secrets
	logg.Debugf("scheduler service config: %+v", app.cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	var storage ConnectableStorage
	if app.cfg.StorageConfig.UseDb {
		storage = sqlstorage.NewStorage(
			app.cfg.StorageConfig.Host,
			app.cfg.StorageConfig.Port,
			app.cfg.StorageConfig.User,
			app.cfg.StorageConfig.Password,
			app.cfg.StorageConfig.Database)
	} else {
		storage = memorystorage.NewStorage()
	}
	if err := storage.Connect(ctx); err != nil {
		logg.Fatalf("can't connect to database %s", err.Error())
	}
	defer storage.Close()

	// events service
	eventsService := events.NewEventsService(storage)

	// init mq
	mqClient := mqclient.NewMQClient(app.cfg.RabbitConfig, logg)
	err = mqClient.Init()
	if err != nil {
		return err
	}
	defer mqClient.Close()

	// cron
	cron := gocron.NewScheduler(time.UTC)

	// cleanup
	cron.Every(time.Duration(app.cfg.TasksConfig.CleanupIntervalMin) * time.Minute).Do(func() {
		logg.Debugf("cleanup task started...")
		_ = eventsService
		until := time.Now().Add(-time.Hour * time.Duration(24*app.cfg.TasksConfig.CleanupPeriodDays))
		if _, err := eventsService.RemoveEventsOlderThan(ctx, until); err != nil && !errors.Is(err, errs.ErrEventNotFound) {
			logg.Errorf("remove old events error %v", err)
		}
		logg.Debugf("cleanup task finished...")
	})

	// notify
	cron.Every(time.Duration(app.cfg.TasksConfig.NotifyCheckIntervalMin) * time.Minute).Do(func() {
		logg.Debugf("notify task started...")
		t := time.Now()
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		notifyTime := t.Add(time.Minute * time.Duration(app.cfg.TasksConfig.NotifyCheckIntervalMin))
		events, err := eventsService.GetEventsForDay(ctx, day)
		if err != nil {
			logg.Errorf("get events for notify error: %v", err)
			return
		}

		for _, e := range events {
			if !e.IsNotified && e.StartDatetime.Compare(notifyTime) <= 0 {
				// marshaling
				msg, err := json.Marshal(e)
				if err != nil {
					logg.Errorf("marshaling error %v", err)
					continue
				}
				// notify
				err = mqClient.Notify(app.cfg.RabbitConfig.Queue, msg)
				if err != nil {
					logg.Errorf("event notify error: %v", err)
					continue
				}
				// mark as notified
				eventsService.SetIsNotified(ctx, e.Id)

				logg.Debugf("event id %d sended to notify queue", e.Id)
			}
		}
		logg.Debugf("notify task finished...")
	})

	cron.StartAsync()

	<-ctx.Done()

	logg.Infof("calendar scheduler stopped")

	return nil
}
