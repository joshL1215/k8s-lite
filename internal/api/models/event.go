package models

type WatchEvent struct {
	EventType   string `json:"eventType"`
	EventObject string `json:"objectType"`
	Pod         *Pod   `json:"pod,omitempty"`
	Node        *Node  `json:"node,omitempty"`
}
