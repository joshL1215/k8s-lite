package models

import "time"

// how enums are done in Go
// Pod phase
type PodPhase string

const (
	PodPending     PodPhase = "Pending"
	PodScheduled   PodPhase = "Scheduled"
	PodRunning     PodPhase = "Running"
	PodTerminating PodPhase = "Terminating"
	PodDeleted     PodPhase = "Deleted"
)

type Pod struct {
	Name              string     `json:"name"`
	Namespace         string     `json:"namespace"`
	Image             string     `json:"image"`
	NodeName          string     `json:"nodeName,omitempty"`
	Phase             PodPhase   `json:"phase"`
	DeletionTimestamp *time.Time `json:"deleteTime,omitempty"`
}
