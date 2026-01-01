package main

import (
	"log"

	"github.com/joshL1215/k8s-lite/internal/api/client"
	"github.com/joshL1215/k8s-lite/internal/api/models"
)

const DefaultNamespace = "default"

var nextNodeIdx = 0

func schedulePod(cl *client.Client, pod *models.Pod) {
	if pod.DeletionTimestamp != nil {
		log.Printf("Scheduler could not schedule pod %s/%s that is marked for deletion", pod.Namespace, pod.Name)
	}

	readyNodes, err := cl.ListNodes(models.NodeReady)
	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		return
	}

	if len(readyNodes) == 0 {
		log.Printf("No ready nodes available to schedule pod %s/%s", pod.Namespace, pod.Name)
		return
	}

	selectedNode := readyNodes[nextNodeIdx%len(readyNodes)]

	updatedPod := *pod
	updatedPod.NodeName = selectedNode.Name
	updatedPod.Phase = models.PodScheduled
	nextNodeIdx++

	if _, err := cl.UpdatePod(&updatedPod); err != nil {
		log.Printf("Error scheduling pod %s/%s to node %s: %v", pod.Namespace, pod.Name, selectedNode.Name, err)
		return
	}
	log.Printf("Scheduled pod %s/%s to node %s", pod.Namespace, pod.Name, selectedNode.Name)
}

func main() {

	log.Print("Starting scheduler...")

	cl, err := client.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
		return
	}

	ch, err := cl.WatchPods(DefaultNamespace)
	if err != nil {
		log.Fatalf("Error watching pods: %v", err)
		return
	}

	log.Print("Scheduler started. Listening for pod events...")

	for event := range ch {
		eventType := event.EventType
		log.Printf("Received event: %v", event)
		if eventType != models.AddEvent {
			log.Print("Ignoring non-add event")
			continue
		}

		pod := event.Pod
		if pod.Phase == models.PodPending {
			schedulePod(cl, pod)
		}
	}
}
