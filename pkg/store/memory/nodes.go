package memory

import (
	"fmt"

	"github.com/joshL1215/k8s-lite/pkg/api"
)

func (s *InMemoryStore) CreateNode(node *api.Node) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.nodes[node.Name]; exists {
		return fmt.Errorf("a node named %s already exists", node.Name)
	}
	s.nodes[node.Name] = node
	return nil
}

func (s *InMemoryStore) GetNode(name string) (*api.Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	node, exists := s.nodes[name]
	if !exists {
		return nil, fmt.Errorf("no node named %s", name)
	}
	return node, nil
}

func (s *InMemoryStore) UpdateNode(node *api.Node) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.nodes[node.Name]; !exists {
		return fmt.Errorf("no node named %s to update", node.Name)
	}
	s.nodes[node.Name] = node
	return nil
}

func (s *InMemoryStore) DeleteNode(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.nodes[name]; !exists {
		return fmt.Errorf("no node named %s to delete", name)
	}
	delete(s.nodes, name)
	return nil
}

func (s *InMemoryStore) ListNodes() ([]*api.Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	nodeList := make([]*api.Node, 0)
	for _, node := range s.nodes {
		nodeList = append(nodeList, node)
	}
	return nodeList, nil
}
