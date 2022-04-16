package email

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/sanches1984/gopkg-common/client/email/provider"
)

type Client struct {
	metricSuccess *prometheus.Counter
	metricFailed  *prometheus.Counter
	sender        provider.ISender
	logger        zerolog.Logger
}

func NewMessage() *provider.Message {
	return &provider.Message{}
}

func NewClient(provider provider.IProvider, fromAddress, fromName string, option ...Option) (*Client, func() error, error) {
	c := Client{}
	for _, o := range option {
		o(&c)
	}
	sender, closer, err := provider.Connect(fromAddress, fromName)
	if err != nil {
		return nil, closer, err
	}
	c.sender = sender
	return &c, closer, nil
}

func (c *Client) Send(ctx context.Context, subject string, addressNameMap map[string]string, msg *provider.Message) error {
	msg.Subject = subject
	msg.To = make(provider.ContactList, 0, len(addressNameMap))
	for address, name := range addressNameMap {
		msg.To = append(msg.To, &provider.Contact{Address: address, Name: name})
	}
	if err := msg.Prepare(ctx); err != nil {
		return err
	}
	err := c.sender.Send(ctx, msg)

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
