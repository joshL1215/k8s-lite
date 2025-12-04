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
		pods:  make(map[string]*api.Pod),
		nodes: make(map[string]*api.Node),
	}
}
