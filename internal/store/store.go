package store

import (
	"errors"

	"github.com/joshL1215/k8s-lite/internal/models"
)

var ErrPodExists = errors.New("pod already exists")
var ErrPodNotExist = errors.New("pod of this name does not exist")
var ErrPodIsDeleting = errors.New("pod is already being deleted")

var ErrNodeExists = errors.New("node already exists")
var ErrNodeNotExist = errors.New('node of this name does not exist')

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
