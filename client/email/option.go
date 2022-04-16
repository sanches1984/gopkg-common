package email

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type Option func(c *Client)

func WithMetricSuccess(metric prometheus.Counter) Option {
	return func(c *Client) {
		c.metricSuccess = &metric
	}
}

func WithMetricFailed(metric prometheus.Counter) Option {
	return func(c *Client) {
		c.metricFailed = &metric
	}
}

func WithLogger(logger zerolog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}
