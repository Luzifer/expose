package ngrok2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL        = "http://localhost:4040"
	defaultRequestTimeout = 5 * time.Second
)

// Client represents the API of the ngrok2 daemon
type Client struct {
	HTTPClient     *http.Client
	BaseURL        string
	RequestTimeout time.Duration
}

// New creates a new Client with default settings
func New() *Client {
	return &Client{
		HTTPClient:     http.DefaultClient,
		BaseURL:        defaultBaseURL,
		RequestTimeout: defaultRequestTimeout,
	}
}

func (c *Client) do(method, subPath string, payload interface{}, response interface{}) error {
	var (
		req    *http.Request
		urlStr = c.BaseURL + subPath
	)

	if payload != nil {
		buf := bytes.NewBufferString("")
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return fmt.Errorf("Unable to encode payload: %s", err)
		}
		req, _ = http.NewRequest(method, urlStr, buf)
	} else {
		req, _ = http.NewRequest(method, urlStr, nil)
	}

	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), c.RequestTimeout)
	defer cancel()

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("Unable to execute request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Received unexpected status code %d on request.", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("Unable to decode response: %s", err)
		}
	}

	return nil
}
