package middleware

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

// NewValidateServerInterceptor wrap endpoint with validation middleware
func NewValidateServerInterceptor(validate *validator.Validate) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = validate.Struct(req)
		if err != nil {
			return
		}

		return handler(ctx, req)
	}
}
