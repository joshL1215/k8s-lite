package store

import (
	"github.com/joshL1215/k8s-lite/pkg/api"
)

// Defines an agnostic store interface
type StoreInterface interface {
	CreatePod(pod *api.Pod) error
	GetPod(namespace, name string) (*api.Pod, error)
	UpdatePod(pod *api.Pod) error
	DeletePod(namespace, name string) error
	ListPods(namespace string) ([]*api.Pod, error)

	CreateNode(node *api.Node) error
	GetNode(name string) (*api.Node, error)
	UpdateNode(node *api.Node) error
	DeleteNode(name string) error
	ListNodes() ([]*api.Node, error)
}
