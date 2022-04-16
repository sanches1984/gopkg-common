package provider

import (
	"context"
)

type IProvider interface {
	Connect() (ISender, func() error, error)
}

type ISender interface {
	Send(ctx context.Context, phone int64, message string) error
}
