package models

type WatchEvent struct {
	Type   string `json:"eventType"`
	Object string `json:"objectType"`
	Pod    *Pod   `json:"pod,omitempty"`
	Node   *Node  `json:"node,omitempty"`
}
