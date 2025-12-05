package models

// Node status enum
type NodeStatus string

const (
	NodeReady    NodeStatus = "Ready"
	NodeNotReady NodeStatus = "Ready"
)

type Node struct {
	Name    string     `json:"name"`
	Address string     `json:"address"`
	Status  NodeStatus `json:"status"`
}
