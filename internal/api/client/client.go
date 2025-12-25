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

// Node operations from client
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

func (c *Client) getNode(nodeName string) (*models.Node, error) {
	req, err := http.NewRequest("GET", c.buildURL("nodes", nodeName), nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request to fetch node: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request to fetch node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch node, status code: %d", resp.StatusCode)
	}

	var fetchedNode models.Node
	if err := json.NewDecoder(resp.Body).Decode(&fetchedNode); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &fetchedNode, nil
}

func (c *Client) listNodes() ([]models.Node, error) {
	req, err := http.NewRequest("GET", c.buildURL("nodes"), nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request to list nodes: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request to list nodes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list nodes, status code: %d", resp.StatusCode)
	}

	var nodes []models.Node
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return nodes, nil
}

func (c *Client) deleteNode(nodeName string) error {
	req, err := http.NewRequest("DELETE", c.buildURL("nodes", nodeName), nil)
	if err != nil {
		return fmt.Errorf("error while creating DELETE request to delete node: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while making DELETE request to delete node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete node, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) updateNode(node *models.Node) (*models.Node, error) {
	body, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling node: %w", err)
	}

	req, err := http.NewRequest("PUT", c.buildURL("nodes", node.Name), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating PUT request to update node: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making PUT request to update node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update node, status code: %d", resp.StatusCode)
	}

	var updatedNode models.Node
	if err := json.NewDecoder(resp.Body).Decode(&updatedNode); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &updatedNode, nil
}

// Pod operations from client
