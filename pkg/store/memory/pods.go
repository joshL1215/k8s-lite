package memory

import (
	"fmt"
	"time"

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
		return fmt.Errorf("pod %s already exists in namespace %s", pod.Name, pod.Namespace)
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
		return nil, fmt.Errorf("no pod with name %s exists in namespace %s", name, namespace)
	}
	return pod, nil
}

// UpdatePod
func (s *InMemoryStore) UpdatePod(pod *api.Pod) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(pod.Namespace, pod.Name)
	currPod, exists := s.pods[key]
	if !exists {
		return fmt.Errorf("no pod with name %s exists in namespace %s", pod.Namespace, pod.Name)
	}

	if currPod.DeletionTimestamp != nil {
		return fmt.Errorf("cannot update pod %s in namespace %s, it is being deleted", pod.Namespace, pod.Name)
	}

	s.pods[key] = pod
	return nil

	// TODO: Needs handling of DeletionTimestamp or terminating pods
	// Needs phase restrictions for terminating pods
	// Needs protection against changing NodeName during termination
	// Needs guidance to use DeletePod for deletion updates
}

// DeletePod
func (s *InMemoryStore) DeletePod(namespace, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(namespace, name)
	currPod, exists := s.pods[key]
	if !exists {
		return fmt.Errorf("no pod with name %s exists in namespace %s", namespace, name)
	}

	if currPod.DeletionTimestamp != nil {
		return fmt.Errorf("cannot delete pod %s in namespace %s, it is already being deleted", namespace, name)
	}
	currTime := time.Now()
	currPod.DeletionTimestamp = &currTime
	currPod.Phase = api.PodTerminating
	s.pods[key] = currPod
	return nil
}

// ListPods
func (s *InMemoryStore) ListPods(namespace string) ([]*api.Pod, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	podList := make([]*api.Pod, 0, len(s.pods))
	for _, pod := range s.pods {
		if pod.Namespace == namespace {
			podList = append(podList, pod)
		}
	}
	return podList, nil
}
