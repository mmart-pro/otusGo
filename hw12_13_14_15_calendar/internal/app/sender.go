package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/mqclient"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/pkg/logger"
)

type Sender struct {
	cfg config.SenderConfig
}

func NewSender(cfg config.SenderConfig) *Sender {
	return &Sender{
		cfg: cfg,
	}
}

func (app Sender) Startup() error {
	// logger
	logg, err := logger.NewLogger(app.cfg.LoggerConfig.Level, app.cfg.LoggerConfig.LogFile)
	if err != nil {
		return fmt.Errorf("can't start logg: %w", err)
	}
	defer logg.Close()

	// IRL exclude secrets
	logg.Debugf("sender service config: %+v", app.cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// init mq
	mqClient := mqclient.NewMQClient(app.cfg.RabbitConfig, logg)
	err = mqClient.Init()
	if err != nil {
		return err
	}
	defer mqClient.Close()

	msgs, err := mqClient.Consume(ctx, app.cfg.RabbitConfig.Queue)
	if err != nil {
		return err
	}
	logg.Debugf("consumer started...")

	for m := range msgs {
		var event model.Event
		if err := json.Unmarshal(m, &event); err != nil {
			logg.Errorf("unmarshaling error %v", err)
			continue
		}
		logg.Infof("notify message: %+v", event)
	}

	<-ctx.Done()

	logg.Infof("calendar sender stopped")

	return nil
}
