// TODO: Build store w/ protobuf later
// Currently in memory key pairing

package memory

import (
	"sync"

	"github.com/joshL1215/k8s-lite/pkg/api"
)

type InMemoryStore struct {
	mutex sync.RWMutex
	pods  map[string]*api.Pod
	nodes map[string]*api.Node
}

func CreateInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		pods: make(map[string]*api.Pod),
		// TODO: nodes later
	}
}

type StoreInterface interface {
	CreatePod(pod *api.Pod) error
	GetPod(namespace, name string) (*api.Pod, error)
	UpdatePod(pod *api.Pod) error
	DeletePod(namespace, name string) error
	ListPods(namespace string) ([]*api.Pod, error)

	// node methods
}
