// TODO: Build store w/ protobuf later
// Currently in memory key pairing

package store

import (
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
