package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	grpcserver "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/grpc"
	httpserver "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/http"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/service/events"
	memorystorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/pkg/logger"
)

type ConnectableStorage interface {
	events.EventsStorage
	Connect(ctx context.Context) error
	Close() error
}

type Calendar struct {
	cfg config.CalendarConfig
}

func NewCalendar(cfg config.CalendarConfig) *Calendar {
	return &Calendar{
		cfg: cfg,
	}
}

func (app Calendar) Startup(ctx context.Context) error {
	// logger
	logg, err := logger.NewLogger(app.cfg.LoggerConfig.Level, app.cfg.LoggerConfig.LogFile)
	if err != nil {
		return fmt.Errorf("can't start logg: %w", err)
	}
	defer logg.Close()

	// IRL exclude secrets
	logg.Debugf("service config: %+v", app.cfg)

	ctx, cancel := signal.NotifyContext(ctx,
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

	logg.Infof("starting calendar api...")

	http := httpserver.NewServer(app.cfg.HttpConfig.GetEndpoint(), app.cfg.GrpcConfig.GetEndpoint(), logg, eventsService)
	go func() {
		if err := http.Start(ctx); err != nil {
			logg.Errorf("failed to start http server %s", err.Error())
			cancel()
		}
	}()

	grpc := grpcserver.NewServer(app.cfg.GrpcConfig.GetEndpoint(), logg, eventsService)
	go func() {
		if err := grpc.Start(); err != nil {
			logg.Errorf("failed to start grpc server %s", err.Error())
			cancel()
		}
	}()

	<-ctx.Done()

	grpc.Stop()

	// 3 sec to stop
	stopContext, stopTimeout := context.WithTimeout(context.Background(), time.Second*3)

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
