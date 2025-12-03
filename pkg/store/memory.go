package store

import (
	"fmt"
	"sync"

	"github.com/joshL1215/k8s-lite/pkg/api"
)

type InMemoryStore struct {
	mutex sync.RWMutex
	pods  map[string]*api.Pod
	// TODO: nodes later
}

func CreateInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		pods: make(map[string]*api.Pod),
		// TODO: nodes later
	}
}

func podKey(pod *api.Pod) string {
	return fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
}

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
