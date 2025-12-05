package models

import "time"

// how enums are done in Go
// Pod phase
type Phase string

const (
	PodPending     Phase = "Pending"
	PodScheduled   Phase = "Scheduled"
	PodRunning     Phase = "Running"
	PodTerminating Phase = "Terminating"
	PodDeleted     Phase = "Deleted"
)

type Pod struct {
	Name              string     `json:"name"`
	Namespace         string     `json:"namespace"`
	Image             string     `json:"image"`
	NodeName          string     `json:"nodeName,omitempty"`
	Phase             Phase      `json:"phase"`
	DeletionTimestamp *time.Time `json:"deleteTime,omitempty"`
}
