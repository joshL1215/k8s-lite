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

func (c *Client) ListPods(namespace string, phase models.PodPhase) ([]models.Pod, error) {
	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	var allPods []models.Pod
	if err := json.NewDecoder(resp.Body).Decode(&allPods); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if phase == "" { // No phase filter, return all
		return allPods, nil
	}

	var filteredPods []models.Pod
	for _, pod := range allPods {
		if pod.Phase == phase {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods, nil
}

// ListNodes fetches nodes, optionally filtering by status.
// Similar to ListPods, filters client-side for simplicity.
func (c *Client) ListNodes(status models.NodeStatus) ([]models.Node, error) {
	urlStr := c.buildURL("api", "v1", "nodes")
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	var allNodes []models.Node
	if err := json.NewDecoder(resp.Body).Decode(&allNodes); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if status == "" { // No status filter, return all
		return allNodes, nil
	}

	var filteredNodes []models.Node
	for _, node := range allNodes {
		if node.Status == status {
			filteredNodes = append(filteredNodes, node)
		}
	}
	return filteredNodes, nil
}

// UpdatePod sends a PUT request to update a pod.
func (c *Client) UpdatePod(pod *models.Pod) error {
	urlStr := c.buildURL("api", "v1", "namespaces", pod.Namespace, "pods", pod.Name)

	body, err := json.Marshal(pod)
	if err != nil {
		return fmt.Errorf("marshalling pod: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO: Read body for more detailed error message from server
		return fmt.Errorf("server returned non-OK status for update: %d", resp.StatusCode)
	}
	// Optionally decode the response body if the updated pod is returned
	return nil
}

// GetNode fetches a specific node by name.
func (c *Client) GetNode(name string) (*models.Node, error) {
	urlStr := c.buildURL("api", "v1", "nodes", name)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for get node: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request for get node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("node %s not found", name) // Specific error for not found
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status for get node: %d", resp.StatusCode)
	}

	var node models.Node
	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		return nil, fmt.Errorf("decoding node response: %w", err)
	}
	return &node, nil
}

// CreatePod sends a POST request to create a pod in a specific namespace.
func (c *Client) CreatePod(namespace string, pod *models.Pod) (*models.Pod, error) {
	if namespace == "" {
		namespace = "default" // Or use a constant
	}
	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods")

	body, err := json.Marshal(pod)
	if err != nil {
		return nil, fmt.Errorf("marshalling pod: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// TODO: Read body for more detailed error message from server
		return nil, fmt.Errorf("server returned non-Created status for create pod: %d", resp.StatusCode)
	}

	var createdPod Pod
	if err := json.NewDecoder(resp.Body).Decode(&createdPod); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &createdPod, nil
}

// GetPod fetches a specific pod by name from a namespace.
func (c *Client) GetPod(namespace, name string) (*models.Pod, error) {
	if namespace == "" {
		namespace = "default"
	}
	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods", name)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for get pod: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request for get pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pod %s/%s not found", namespace, name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status for get pod: %d", resp.StatusCode)
	}

	var pod Pod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, fmt.Errorf("decoding pod response: %w", err)
	}
	return &pod, nil
}

// DeletePod sends a DELETE request to remove a pod.
func (c *Client) DeletePod(namespace, name string) error {
	if namespace == "" {
		namespace = "default"
	}
	urlStr := c.buildURL("api", "v1", "namespaces", namespace, "pods", name)

	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
	if err != nil {
		return fmt.Errorf("creating request for delete pod: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request for delete pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent { // Some APIs return 204 for delete
		// TODO: Read body for more detailed error message from server
		return fmt.Errorf("server returned non-OK status for delete pod: %d", resp.StatusCode)
	}
	return nil
}
