package store

import "github.com/joshL1215/k8s-lite/internal/models"

// Defines an agnostic store interface
type StoreInterface interface {
	CreatePod(pod *models.Pod) error
	GetPod(namespace, name string) (*models.Pod, error)
	UpdatePod(pod *models.Pod) error
	DeletePod(namespace, name string) error
	ListPods(namespace string) ([]*models.Pod, error)

	CreateNode(node *models.Node) error
	GetNode(name string) (*models.Node, error)
	UpdateNode(node *models.Node) error
	DeleteNode(name string) error
	ListNodes() ([]*models.Node, error)
}
