package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/joshL1215/k8s-lite/internal/api/models"
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

func (c *Client) createNode(node *models.Node) (*models.Node, error) {
	body, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling node: %w", err)
	}

	req, err := http.NewRequest("POST", c.buildURL("nodes"), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating POST request to register node: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making POST request to register node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to register node, status code: %d", resp.StatusCode)
	}

	var createdNode models.Node
	if err := json.NewDecoder(resp.Body).Decode(&createdNode); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &createdNode, nil
}
