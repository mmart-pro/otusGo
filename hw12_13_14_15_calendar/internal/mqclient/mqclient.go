package mqclient

import (
	"context"
	"fmt"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type MQClient struct {
	cfg    config.RabbitConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger Logger
}

func NewMQClient(cfg config.RabbitConfig, logger Logger) MQClient {
	return MQClient{
		cfg:    cfg,
		logger: logger,
	}
}

func (c *MQClient) Init() error {
	var err error
	c.conn, err = amqp.Dial("amqp://" + c.cfg.User + ":" + c.cfg.Password + "@" + c.cfg.Host + ":" + c.cfg.Port + "/")
	if err != nil {
		return err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return err
	}

	err = c.ch.ExchangeDeclare(
		c.cfg.Queue, // name
		"fanout",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(
		c.cfg.Queue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(
		c.cfg.Queue, // queue name
		"",          // routing key
		c.cfg.Queue, // exchange
		false,       // nowait
		nil,         // args
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *MQClient) Notify(exchange string, msg []byte) error {
	return c.ch.Publish(
		exchange, // exchange
		"",       // key
		true,     // mandatory
		false,    // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json; charset=utf-8",
			Body:         msg,
		})
}

func (c *MQClient) Consume(ctx context.Context, queue string) (<-chan []byte, error) {
	result := make(chan []byte)

	go func() {
		<-ctx.Done()
		c.ch.Close()
	}()

	incoming, err := c.ch.Consume(queue, queue, false, false, false, false, nil)
	if err != nil {
		e := fmt.Errorf("start consuming error: %w", err)
		c.logger.Errorf(e.Error())
		return nil, e
	}

	go func() {
		defer func() {
			close(result)
			c.logger.Debugf("result channel closed")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-incoming:
				if err := ev.Ack(false); err != nil {
					c.logger.Errorf("ack error %v", err)
					continue
				}

				select {
				case <-ctx.Done():
					return
				case result <- ev.Body:
				}
			}
		}
	}()

	return result, nil
}

func (c *MQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}
