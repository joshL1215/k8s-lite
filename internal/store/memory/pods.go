package memory

import (
	"fmt"
	"time"

	"github.com/joshL1215/k8s-lite/internal/api/models"
	"github.com/joshL1215/k8s-lite/internal/store"
)

func podKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

// CreatePod
func (s *InMemoryStore) CreatePod(pod *models.Pod) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(pod.Namespace, pod.Name)
	if _, exists := s.pods[key]; exists {
		return fmt.Errorf("%w: pod %s already exists in namespace %s", store.ErrPodExists, pod.Name, pod.Namespace)
	}
	s.pods[key] = pod
	return nil
}

// GetPod
func (s *InMemoryStore) GetPod(namespace, name string) (*models.Pod, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	key := podKey(namespace, name)
	pod, exists := s.pods[key]
	if !exists {
		return nil, fmt.Errorf("%w :no pod with name %s exists in namespace %s", store.ErrPodNotExist, name, namespace)
	}
	return pod, nil
}

// UpdatePod
func (s *InMemoryStore) UpdatePod(pod *models.Pod) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(pod.Namespace, pod.Name)
	currPod, exists := s.pods[key]
	if !exists {
		return fmt.Errorf("%w: no pod with name %s exists in namespace %s", store.ErrPodNotExist, pod.Name, pod.Namespace)
	}

	if currPod.DeletionTimestamp != nil {
		return fmt.Errorf("%w: cannot update pod %s in namespace %s, it is being deleted", store.ErrPodIsDeleting, pod.Namespace, pod.Name)
	}

	s.pods[key] = pod
	return nil
}

// DeletePod
func (s *InMemoryStore) DeletePod(namespace, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := podKey(namespace, name)
	currPod, exists := s.pods[key]
	if !exists {
		return fmt.Errorf("%w: no pod with name %s exists in namespace %s", store.ErrPodNotExist, namespace, name)
	}

	if currPod.DeletionTimestamp != nil {
		return fmt.Errorf("%w: cannot delete pod %s in namespace %s, it is already being deleted", store.ErrPodIsDeleting, namespace, name)
	}
	currTime := time.Now()
	currPod.DeletionTimestamp = &currTime
	currPod.Phase = models.PodTerminating
	s.pods[key] = currPod
	return nil
}

// ListPods
func (s *InMemoryStore) ListPods(namespace string) ([]*models.Pod, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	podList := make([]*models.Pod, 0, len(s.pods))
	for _, pod := range s.pods {
		if pod.Namespace == namespace {
			podList = append(podList, pod)
		}
	}
	return podList, nil
}
