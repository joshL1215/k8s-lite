package store

import (
	"fmt"

	"github.com/joshL1215/k8s-lite/pkg/api"
)

// CreatePod
func (s *InMemoryStore) CreatePod(pod *api.Pod) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(pod)
	if _, exists := s.pods[key]; exists {
		return fmt.Errorf("Pod %s already exists in namespace %s", pod.Name, pod.Namespace)
	}
	s.pods[key] = pod
	return nil
}

// GetPod

// UpdatePod

// DeletePod

// ListPods

//
