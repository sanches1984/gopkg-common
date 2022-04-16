package sms

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/gopkg-common/client/sms/provider"

	"github.com/prometheus/client_golang/prometheus"
)

type Client struct {
	metricSuccess *prometheus.Counter
	metricFailed  *prometheus.Counter
	sender        provider.ISender
	logger        zerolog.Logger
}

func NewClient(provider provider.IProvider, option ...Option) (*Client, func() error, error) {
	c := &Client{}
	for _, o := range option {
		o(c)
	}
	sender, closer, err := provider.Connect()
	if err != nil {
		return nil, closer, err
	}
	c.sender = sender
	return c, closer, nil
}

func (c Client) Send(ctx context.Context, phone int64, message string) error {
	err := c.sender.Send(ctx, phone, message)
	if err == nil {
		if c.metricSuccess != nil {
			(*c.metricSuccess).Inc()
		}
	} else {
		if c.metricFailed != nil {
			(*c.metricFailed).Inc()
		}
	}
	return err
}
