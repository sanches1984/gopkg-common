// +build payment

package sms

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sanches1984/gopkg-common/client/sms/provider"
	"github.com/stretchr/testify/assert"
)

func TestMqClient(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env variables, error: " + err.Error())
	}

	client, _, err := NewClient(provider.NewMqProvider(
		os.Getenv("SMS_MQ_DSN"),
		os.Getenv("SMS_MQ_EXCHANGE"),
		os.Getenv("SMS_MQ_ROUTING_KEY"),
		zerolog.Nop(),
	))
	assert.Nil(t, err)
	phone, err := strconv.ParseInt(os.Getenv("SMS_PHONE"), 10, 64)
	assert.Nil(t, err)
	err = client.Send(context.Background(), phone, "test")
	assert.Nil(t, err)
}

func TestTerraClient(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env variables, error: " + err.Error())
	}

	client, _, err := NewClient(provider.NewTerraProvider(
		os.Getenv("SMS_TERRA_URL"),
		os.Getenv("SMS_TERRA_LOGIN"),
		os.Getenv("SMS_TERRA_SENDER"),
		os.Getenv("SMS_TERRA_PASSWORD"),
		zerolog.Nop(),
	))
	assert.Nil(t, err)
	phone, err := strconv.ParseInt(os.Getenv("SMS_PHONE"), 10, 64)
	assert.Nil(t, err)
	err = client.Send(context.Background(), phone, "test")
	assert.Nil(t, err)
}
