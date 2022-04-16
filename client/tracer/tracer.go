package tracer

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"io"
)

var appTracer *opentracing.Tracer

func InitTracer(serviceName, addr string) (io.Closer, error) {
	transport, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, err
	}
	tracer, tracerCloser := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
	)
	appTracer = &tracer
	opentracing.InitGlobalTracer(tracer)
	return tracerCloser, nil
}

func GetTracer() opentracing.Tracer {
	return *appTracer
}

func StartSpan(ctx context.Context, operationName string) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span = (*appTracer).StartSpan(operationName, opentracing.ChildOf(span.Context()))
	} else {
		span = (*appTracer).StartSpan(operationName)
	}

	return span
}
