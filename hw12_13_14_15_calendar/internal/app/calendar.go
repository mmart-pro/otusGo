package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/logger"
	httpserver "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/http"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/service/events"
	memorystorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/sql"
)

type ConnectableStorage interface {
	events.EventsStorage
	Connect(ctx context.Context) error
	Close() error
}

type Calendar struct {
	configFile string
}

func NewCalendar(configFile string) *Calendar {
	return &Calendar{
		configFile: configFile,
	}
}

func (app Calendar) Startup(ctx context.Context) error {
	// config
	cfg, err := config.NewConfig(app.configFile)
	if err != nil {
		return fmt.Errorf("error read config: %s\n", err)
	}

	// logger
	logg, err := logger.NewLogger(cfg.LoggerConfig.Level, cfg.LoggerConfig.LogFile)
	if err != nil {
		return fmt.Errorf("can't start logg: %s\n", err)
	}
	defer logg.Close()

	// IRL exclude secrets
	logg.Debugf("service config: %+v", cfg)

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	var storage ConnectableStorage
	if cfg.StorageConfig.UseDb {
		storage = sqlstorage.NewStorage(
			cfg.StorageConfig.Host,
			cfg.StorageConfig.Port,
			cfg.StorageConfig.User,
			cfg.StorageConfig.Password,
			cfg.StorageConfig.Database)
	} else {
		storage = memorystorage.NewStorage()
	}
	if err := storage.Connect(ctx); err != nil {
		logg.Fatalf("can't connect to database %s", err.Error())
	}
	defer storage.Close()

	// events service
	eventsService := events.NewEventsService(logg, storage)

	logg.Infof("starting calendar api...")

	http := httpserver.NewServer(cfg.HttpConfig, logg, eventsService)
	go func() {
		if err := http.Start(ctx); err != nil {
			logg.Errorf("failed to start http server %s", err.Error())
			cancel()
		}
	}()

	<-ctx.Done()

	// 3 sec to stop
	stopContext, stopTimeout := context.WithTimeout(context.Background(), time.Second*3)
	defer stopTimeout()

	logg.Debugf("stopping calendar api...")

	go func() {
		defer stopTimeout()
		if err := http.Stop(stopContext); err != nil {
			logg.Errorf("failed to stop http server: " + err.Error())
		}
	}()

	<-stopContext.Done()
	logg.Infof("calendar api stopped")

	return nil
}
