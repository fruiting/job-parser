package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't make get request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read all from body: %w", err)
	}

	return body, nil
}

func (c *Client) Post(url string, body []byte) error {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unexpected status: %s", resp.Status))
	}

	return nil
}
