package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func (c *Client) ShowURL() string {
	return c.baseURL.String()
}

// Node operations from client
func (c *Client) CreateNode(node *models.Node) (*models.Node, error) {
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

func (c *Client) GetNode(nodeName string) (*models.Node, error) {
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

func (c *Client) ListNodes() ([]models.Node, error) {
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

func (c *Client) DeleteNode(nodeName string) error {
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

func (c *Client) UpdateNode(node *models.Node) (*models.Node, error) {
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

func (c *Client) CreatePod(pod *models.Pod) (*models.Pod, error) {
	body, err := json.Marshal(pod)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling pod: %w", err)
	}

	urlStr := c.buildURL("api", "v1", "namespaces", pod.Namespace, "pods")
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating POST request to create pod: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making POST request to create pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create pod, status code: %d", resp.StatusCode)
	}

	var createdPod models.Pod
	if err := json.NewDecoder(resp.Body).Decode(&createdPod); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &createdPod, nil
}

func (c *Client) GetPod(namespace, podName string) (*models.Pod, error) {
	if namespace == "" {
		namespace = "default"
	}

	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods", podName)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request to fetch pod: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request to fetch pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pod, status code: %d", resp.StatusCode)
	}

	var fetchedPod models.Pod
	if err := json.NewDecoder(resp.Body).Decode(&fetchedPod); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &fetchedPod, nil
}

func (c *Client) ListPods(namespace string) ([]models.Pod, error) {
	if namespace == "" {
		namespace = "default"
	}

	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods")
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request to list pods: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request to list pods: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list pods, status code: %d", resp.StatusCode)
	}

	var pods []models.Pod
	if err := json.NewDecoder(resp.Body).Decode(&pods); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return pods, nil
}

func (c *Client) DeletePod(namespace, podName string) error {
	if namespace == "" {
		namespace = "default"
	}

	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods", podName)
	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return fmt.Errorf("error while creating DELETE request to delete pod: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while making DELETE request to delete pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete pod, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) UpdatePod(pod *models.Pod) (*models.Pod, error) {
	body, err := json.Marshal(pod)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling pod: %w", err)
	}

	urlStr := c.buildURL("api", "v1", "namespaces", pod.Namespace, "pods", pod.Name)
	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating PUT request to update pod: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making PUT request to update pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update pod, status code: %d", resp.StatusCode)
	}

	var updatedPod models.Pod
	if err := json.NewDecoder(resp.Body).Decode(&updatedPod); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %w", err)
	}
	return &updatedPod, nil
}

func (c *Client) WatchPods(namespace string) (<-chan models.WatchEvent, error) {
	if namespace == "" {
		namespace = "default"
	}

	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods?watch=true")
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request to watch pods: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request to watch pods: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to watch pods, status code: %d", resp.StatusCode)
	}

	events := make(chan models.WatchEvent)
	go func() {
		defer resp.Body.Close()
		defer close(events)

		decoder := json.NewDecoder(resp.Body)

		for {
			var event models.WatchEvent
			if err := decoder.Decode(&event); err != nil {
				log.Printf("Error decoding watch event: %v", err)
				return
			}
			if event.Object == "pod" {
				events <- event
			}
		}
	}()
	return events, nil
}
