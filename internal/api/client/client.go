package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(urlStr string) (*Client, error) {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error while parsing base URL: %w", err.Error())
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Client) buildURL(segments ...string) string {
	path := c.baseURL.Path
	for _, segment := range segments {
		path = fmt.Sprintf("%s/%s", path, segment)
	}
	// necessary to avoid altering the actual url
	newURL := *c.baseURL
	newURL.Path = path
	return newURL.String()
}
