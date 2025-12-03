package store

import (
	"fmt"

	"github.com/joshL1215/k8s-lite/pkg/api"
)

func podKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

// CreatePod
func (s *InMemoryStore) CreatePod(pod *api.Pod) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(pod.Namespace, pod.Name)
	if _, exists := s.pods[key]; exists {
		return fmt.Errorf("Pod %s already exists in namespace %s", pod.Name, pod.Namespace)
	}
	s.pods[key] = pod
	return nil
}

// GetPod
func (s *InMemoryStore) GetPod(namespace, name string) (*api.Pod, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(namespace, name)
	pod, exists := s.pods[key]
	if !exists {
		return nil, fmt.Errorf("No pod with name %s exists in namespace %s", name, namespace)
	}
	return pod, nil
}

// UpdatePod

// DeletePod

// ListPods

//
