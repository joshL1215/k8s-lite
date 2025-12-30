package models

type EventType string
type EventObject string

const (
	AddEvent          EventType = "ADDED"
	ModificationEvent EventType = "MODIFIED"
	DeletionEvent     EventType = "DELETED"
)

type WatchEvent struct {
	EventType   EventType   `json:"eventType"`
	EventObject EventObject `json:"objectType"`
	Pod         *Pod        `json:"pod,omitempty"`
	Node        *Node       `json:"node,omitempty"`
}
