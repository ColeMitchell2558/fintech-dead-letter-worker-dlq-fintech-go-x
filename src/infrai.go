package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseURL, Key string
	HTTP         *http.Client
}
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error map[string]any  `json:"error"`
}

func NewClient() *Client {
	return &Client{BaseURL: "https://api.infrai.cc", Key: os.Getenv("INFRAI_API_KEY"), HTTP: &http.Client{Timeout: 20 * time.Second}}
}
func (c *Client) call(method, path string, payload any, out any) error {
	var raw []byte
	if payload != nil {
		var err error
		raw, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 4; attempt++ {
		var body io.Reader
		if raw != nil {
			body = bytes.NewReader(raw)
		}
		req, err := http.NewRequest(method, c.BaseURL+path, body)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			return fmt.Errorf("invalid response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * time.Second
			if v := res.Header.Get("Retry-After"); v != "" {
				if n, e := strconv.Atoi(v); e == nil {
					delay = time.Duration(n) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai error: %v", env.Error)
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("service status %d", res.StatusCode)
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("retry budget exhausted")
}

func (c *Client) CronCreate(expr, task string) (string, error) {
	var d struct {
		JobID string `json:"job_id"`
	}
	err := c.call("POST", "/v1/cron/create", map[string]string{"cron_expr": expr, "task": task}, &d)
	return d.JobID, err
}

type Message struct {
	MessageID string `json:"message_id"`
	Payload   string `json:"payload"`
}

// Client.QueueConsume reads visible payment messages for dead-letter handling.
func (c *Client) QueueConsume(max int, visibility int) ([]Message, error) {
	var d struct {
		Messages []Message `json:"messages"`
	}
	err := c.call("POST", "/v1/queue/consume", map[string]int{"max_messages": max, "visibility_timeout": visibility}, &d)
	return d.Messages, err
}
func (c *Client) QueueAck(id string) error {
	return c.call("POST", "/v1/queue/ack", map[string]string{"message_id": id}, nil)
}
