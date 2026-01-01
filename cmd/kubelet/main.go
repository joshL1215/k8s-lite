package main

import (
	"flag"
	"log"
	"time"

	"github.com/joshL1215/k8s-lite/internal/kubelet"
)

// Kubelet will reconcile pod state for its specific node on the interval as well as on node-associated pod events
const syncInterval = 10 * time.Second

const DefaultNamespace = "default"

func main() {
	nodeName := flag.String("node-name", "", "Name of the node being registered")
	nodeAddress := flag.String("node-address", "http://localhost:8081", "Address of the node being registered")
	apiAddress := flag.String("api-server-url", "http://localhost:8080", "URL of the API server")
	flag.Parse()

	if *nodeName == "" {
		log.Fatalf("-node-name flag is required")
	}

	log.Printf("Kubelet starting for node %s at %node address %s, API server at %s", *nodeName, *nodeAddress, *apiAddress)

	k, err := kubelet.NewKubelet(*nodeName, *nodeAddress, *apiAddress)
	if err != nil {
		log.Fatalf("Error creating kubelet: %v", err)
	}

	if err := k.RegisterNode(); err != nil {
		log.Fatalf("Error registering node: %v", err)
	}

	log.Printf("Successfully registed node %s. Kubelet will synchronize pod state on schedule events and on interval of %v", *nodeName, syncInterval)

	log.Printf("Attempting to watch pods for scheduling events...")
	ch, err := k.Client.WatchPods(DefaultNamespace)
	if err != nil {
		log.Fatalf("Error watching pods: %v", err)
	}

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case event := <-ch:
			eventPod := *event.Pod
			log.Printf("Received event: %v", event)

			if eventPod.NodeName != *nodeName {
				log.Printf("Ignoring pod event not associated with registered node")
				continue
			}

			log.Printf("Detected a pod event on registered node, synchronizing pod state")
			k.SyncPods()

		case <-ticker.C:
			log.Printf("Periodic sync interval reached, synchronizing pod state")
			k.SyncPods()
		}
	}

}
