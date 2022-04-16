package provider

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	errors "github.com/sanches1984/gopkg-errors"
	"github.com/streadway/amqp"
	uuid "gopkg.in/satori/go.uuid.v1"
)

type mqProvider struct {
	dsn        string
	channel    *amqp.Channel
	from       string
	exchange   string
	routingKey string
	logger     zerolog.Logger
}

type mqData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func NewMqProvider(dsn, exchange, routingKey string, logger zerolog.Logger) IProvider {
	return &mqProvider{
		dsn:        dsn,
		exchange:   exchange,
		routingKey: routingKey,
		logger:     logger,
	}
}

func (c mqProvider) Connect(fromAddress, fromName string) (ISender, func() error, error) {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return nil, func() error { return nil },
			errors.Internal.ErrWrap(context.Background(), "email: error on AMQP connect", err).WithLogKV("dsn", c.dsn)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, func() error {
				c.logger.Debug().Msg("email: close AMQP connect")
				return conn.Close()
			},
			errors.Internal.ErrWrap(context.Background(), "email: error on amqp channel connect", err).WithLogKV("dsn", c.dsn)
	}

	c.logger.Debug().Str("dsn", c.dsn).Msg("email: use AMPQ connect")

	c.channel = channel
	c.from = Contact{Address: fromAddress, Name: fromName}.String()
	closer := func() error {
		c.logger.Debug().Msg("email: close AMQP connect")
		if err := conn.Close(); err != nil {
			return err
		}
		return channel.Close()
	}
	return c, closer, nil
}

func (c mqProvider) Send(ctx context.Context, msg *Message) error {
	body, err := json.Marshal(mqData{
		From:    c.from,
		To:      msg.To.String(),
		Subject: msg.Subject,
		Content: msg.bodyHTML,
	})
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
			Msg("email: error on publish message")
		return err
	}

	c.logger.Info().
		Str("key", c.routingKey).
		Str("body", string(body)).
		Str("exchange", c.exchange).
		Msg("email: success publish message")

	return nil
}
