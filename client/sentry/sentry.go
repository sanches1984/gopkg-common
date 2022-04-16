package sentry

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	errors "github.com/sanches1984/gopkg-errors"
	"time"
)

var instance *sentryConfig

type sentryConfig struct {
	flushTimeout time.Duration
}

func Init(project, token string, flushTimeout time.Duration) error {
	dsn := fmt.Sprintf("https://%s@sentry.io/%s", token, project)
	err := sentry.Init(sentry.ClientOptions{Dsn: dsn})
	if err != nil {
		return errors.Internal.ErrWrap(context.Background(), "Sentry initialization failed", err).
			WithLogKV("project", project, "token", token, "dsn", dsn)
	}

	instance = &sentryConfig{flushTimeout: flushTimeout}
	return nil
}

func ShouldBeProcessed(err error) bool {
	return errors.IsInternal(err)
}

func Error(err error, tagKV ...string) {
	if instance == nil || sentry.CurrentHub() == nil {
		return
	}
	if ShouldBeProcessed(err) {
		ConfigureScope(tagKV...)
		sentry.CaptureException(err)
		sentry.Flush(instance.flushTimeout)
	}
}

func ConfigureScope(tagKV ...string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		for i := 0; i < len(tagKV); i += 2 {
			scope.SetTag(tagKV[i], tagKV[i+1])
		}
	})
}

func Panic(err interface{}, tagKV ...string) {
	if instance == nil {
		return
	}
	if hub := sentry.CurrentHub(); hub != nil {
		ConfigureScope(tagKV...)
		hub.Recover(fmt.Sprintf("%#v", err))
		sentry.Flush(instance.flushTimeout)
	}
}
