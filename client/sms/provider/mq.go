package provider

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strconv"

	errors "github.com/sanches1984/gopkg-errors"
	"github.com/streadway/amqp"
)

type mqProvider struct {
	channel    *amqp.Channel
	dsn        string
	exchange   string
	routingKey string
	logger     zerolog.Logger
}

type mqData struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func NewMqProvider(dsn, exchange, routingKey string, logger zerolog.Logger) IProvider {
	return &mqProvider{
		dsn:        dsn,
		exchange:   exchange,
		routingKey: routingKey,
		logger:     logger,
	}
}

func (c mqProvider) Connect() (ISender, func() error, error) {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return nil, func() error {
				return nil
			},
			errors.Internal.ErrWrap(context.Background(), "sms: error on AMQP connect", err).WithLogKV("dsn", c.dsn)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, func() error {
				c.logger.Debug().Msg("sms: close AMQP connect")
				return conn.Close()
			},
			errors.Internal.ErrWrap(context.Background(), "sms: error on amqp channel connect", err).WithLogKV("dsn", c.dsn)
	}

	c.logger.Debug().Str("dsn", c.dsn).Msg("sms: use AMPQ connect %v")

	c.channel = channel
	closer := func() error {
		c.logger.Debug().Msg("sms: close AMQP connect")
		if err := conn.Close(); err != nil {
			return err
		}
		return channel.Close()
	}
	return c, closer, nil
}

func (c mqProvider) Send(ctx context.Context, phone int64, message string) error {
	body, err := json.Marshal(mqData{Phone: strconv.FormatInt(phone, 10), Message: message})
	if err != nil {
		return err
	}
	err = c.channel.Publish(
		c.exchange,
		c.routingKey,
		true,
		false,
		amqp.Publishing{
			MessageId:   uuid.NewV4().String(),
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		c.logger.Error().Err(err).
			Str("key", c.routingKey).
			Str("body", string(body)).
			Str("exchange", c.exchange).
			Msg("sms: error on publish message")
		return err
	}

	c.logger.Info().
		Str("key", c.routingKey).
		Str("body", string(body)).
		Str("exchange", c.exchange).
		Msg("sms: success publish message")

	return nil
}
