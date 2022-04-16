package provider

import (
	"context"
)

type IProvider interface {
	Connect(fromAddress, fromName string) (ISender, func() error, error)
}

type ISender interface {
	Send(ctx context.Context, msg *Message) error
}
