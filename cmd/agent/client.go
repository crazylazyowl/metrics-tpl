package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/security"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type client struct {
	baseURL    string
	retries    int
	secret     []byte
	firstDelay time.Duration
}

type clientOptions struct {
	BaseURL string
	Secret  string
}

func newClient(opts clientOptions) *client {
	return &client{
		baseURL:    opts.BaseURL,
		retries:    4,
		firstDelay: time.Second * 1,
		secret:     []byte(opts.Secret),
	}
}

// NOTE: hack for second test
func (c *client) Hack(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/update/", c.baseURL)
	resp, err := http.Head(url)
	if err != nil {
		c.retries = 1
		c.firstDelay = 0
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *client) SendOne(ctx context.Context, m metrics.Metric) error {
	if err := m.Validate(); err != nil {
		return err
	}
	data, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	return c.send(ctx, "/update/", false, data)
}

func (c *client) SendMany(ctx context.Context, mm []metrics.Metric) error {
	for _, m := range mm {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	data, err := json.Marshal(&mm)
	if err != nil {
		return err
	}
	return c.send(ctx, "/updates/", true, data)
}

func (c *client) send(ctx context.Context, endpoint string, compress bool, data []byte) error {
	if compress {
		buf := bytes.NewBuffer(nil)
		w := gzip.NewWriter(buf)
		w.Write(data)
		w.Close()
		data = buf.Bytes()
	}
	var err error
	delay := c.firstDelay
	url := fmt.Sprintf("http://%s%s", c.baseURL, endpoint)
	for range c.retries {
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to prepare request; %w", err)
		}
		var resp *http.Response
		req.Header.Add("Content-Type", "application/json")
		if compress {
			req.Header.Add("Content-Encoding", "gzip")
		}
		if len(c.secret) != 0 {
			req.Header.Add("HashSHA256", security.HMACString([]byte(c.secret), data))
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("retry report")
			time.Sleep(delay)
			delay += time.Second * 2
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("status code = %d, data = %s", resp.StatusCode, string(body))
		}
		break
	}
	return err
}
