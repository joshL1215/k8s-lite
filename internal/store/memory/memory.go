package memory

import (
	"sync"

	"github.com/joshL1215/k8s-lite/internal/api/models"
)

type InMemoryStore struct {
	mutex sync.RWMutex
	pods  map[string]*models.Pod
	nodes map[string]*models.Node
}

func CreateInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		pods:  make(map[string]*models.Pod),
		nodes: make(map[string]*models.Node),
	}
}
