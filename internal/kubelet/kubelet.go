package kubelet

import (
	"fmt"
	"log"

	"github.com/joshL1215/k8s-lite/internal/api/client"
	"github.com/joshL1215/k8s-lite/internal/api/models"
	"github.com/joshL1215/k8s-lite/internal/store"
)

const DefaultNamespace = "default"

type Kubelet struct {
	NodeName    string
	NodeAddress string
	Client      *client.Client
}

func NewKubelet(nodeName, nodeAddress, apiURL string) (*Kubelet, error) {
	cl, err := client.NewClient(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return &Kubelet{
		NodeName:    nodeName,
		NodeAddress: nodeAddress,
		Client:      cl,
	}, nil
}

func (k *Kubelet) RegisterNode() error {
	node := &models.Node{
		Name:    k.NodeName,
		Address: k.NodeAddress,
		Status:  models.NodeReady,
	}
	registeredNode, err := k.Client.CreateNode(node)
	if err == store.ErrNodeExists {
		log.Printf("Node %s already exists, attempting to update...: %v", k.NodeName, err)

		registeredNode, err = k.Client.UpdateNode(node)
		if err != nil {
			return fmt.Errorf("failed to update existing node %s: %v", k.NodeName, err)
		}
		log.Printf("Node %s updated successfully", k.NodeName)
		return nil
	}
	log.Printf("Node %s successfully registered", registeredNode.Name)
	return nil
}

func (k *Kubelet) SyncPods() {
	allPods, err := k.Client.ListPods(DefaultNamespace, "")
	if err != nil {
		log.Printf("Error listing pods from API server: %v", err)
		return
	}

	for _, pod := range allPods {

		if pod.NodeName != k.NodeName {
			continue
		}
		updatingPod := pod

		switch pod.Phase {
		case models.PodTerminating:

			if updatingPod.DeletionTimestamp != nil {

				log.Printf("Pod %s/%s is terminating. Deleting pod...", updatingPod.Namespace, updatingPod.Name)
				updatingPod.Phase = models.PodDeleted

				if _, err := k.Client.UpdatePod(&updatingPod); err != nil {
					log.Printf("Error updating pod %s/%s to Deleted: %v", updatingPod.Namespace, updatingPod.Name, err)
				} else {
					log.Printf("Successfully updated pod %s/%s to Deleted", updatingPod.Namespace, updatingPod.Name)
				}
			}

		case models.PodScheduled:
			log.Printf("Pod %s/%s is scheduled on this node. Starting pod...", updatingPod.Namespace, updatingPod.Name)
			updatingPod.Phase = models.PodRunning

			if _, err := k.Client.UpdatePod(&updatingPod); err != nil {
				log.Printf("Error updating pod %s/%s to Running: %v", updatingPod.Namespace, updatingPod.Name, err)
			} else {
				log.Printf("Successfully updated pod %s/%s to Running", updatingPod.Namespace, updatingPod.Name)
			}

		default:
			log.Printf("Pod %s/%s is in phase %s. No action taken.", updatingPod.Namespace, updatingPod.Name, updatingPod.Phase)
		}
	}
}
