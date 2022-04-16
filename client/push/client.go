package push

import (
	"bytes"
	"context"
	"crypto/md5"
	jsn "encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"net/http"
)

var errTokenExpired = errors.New("NotRegistered")

type Client struct {
	url           string
	apiKey        string
	senderID      string
	metricSuccess *prometheus.Counter
	metricFailed  *prometheus.Counter
	logger        zerolog.Logger
}

type Message struct {
	Token        string `json:"to"`
	Notification struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		Icon  string `json:"icon"`
	} `json:"notification"`
	Data struct {
		Action string `json:"action"`
	} `json:"data"`
}

type pushResponse struct {
	Success int          `json:"success"`
	Failure int          `json:"failure"`
	Results []pushResult `json:"results"`
}

type pushResult struct {
	Error string `json:"error"`
}

func NewClient(url, apiKey, senderID string, option ...Option) *Client {
	c := &Client{
		url:      url,
		apiKey:   apiKey,
		senderID: senderID,
	}
	for _, o := range option {
		o(c)
	}
	return c
}

func (p Message) getMessage() []byte {
	b, _ := jsn.Marshal(p)
	return b
}

func (p Message) getHash() string {
	h := md5.New()
	_, _ = io.WriteString(h, fmt.Sprintf("%v", p))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c Client) Send(ctx context.Context, msg Message) error {
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(msg.getMessage()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+c.apiKey)
	req.Header.Set("Sender", "id="+c.senderID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.logger.Debug().Str("hash", msg.getHash()).Interface("msg", msg).Msg("push request")
	body, _ := ioutil.ReadAll(resp.Body)
	c.logger.Debug().Str("hash", msg.getHash()).Str("body", string(body)).Msg("push request")
	err = c.isSuccess(ctx, body)

	if err == nil {
		if c.metricSuccess != nil {
			(*c.metricSuccess).Inc()
		}
	} else {
		if c.metricFailed != nil {
			(*c.metricFailed).Inc()
		}
	}

	return nil
}

func IsTokenExpired(err error) bool {
	return err == errTokenExpired
}

func (c Client) isSuccess(ctx context.Context, data []byte) error {
	var resp pushResponse
	err := jsn.Unmarshal(data, &resp)
	if err != nil {
		return errors.New("Unable to read push response: " + string(data))
	}

	if resp.Success != 1 {
		if resp.Failure == 1 && resp.isTokenExpired() {
			return errTokenExpired
		}
		return fmt.Errorf("Unable to send push: %v", resp.Results)
	}

	return nil
}

func (resp pushResponse) isTokenExpired() bool {
	for _, err := range resp.Results {
		if err.Error == errTokenExpired.Error() {
			return true
		}
	}

	return false
}
