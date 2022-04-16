// +build payment

package push

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPushSend(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env variables, error: " + err.Error())
	}

	client := NewClient(os.Getenv("PUSH_URL"), os.Getenv("PUSH_API_KEY"), os.Getenv("PUSH_SENDER_ID"))

	msg := Message{Token: os.Getenv("PUSH_TOKEN")}
	msg.Notification.Title = "Hello"
	msg.Notification.Text = "Big big text"
	err = client.Send(context.Background(), msg)
	assert.Nil(t, err)
}
