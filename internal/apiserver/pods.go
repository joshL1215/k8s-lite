package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshL1215/k8s-lite/internal/api/models"
	"github.com/joshL1215/k8s-lite/internal/store"
)

func (s *APIServer) createPodHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var pod models.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "detail": err.Error()})
		return
	}

	if pod.Name == "" {
		c.JSON(400, gin.H{"error": "A pod name must be provided"})
		return
	}
	if namespace == "" {
		namespace = DefaultNamespace
	}
	pod.Namespace = namespace
	pod.Phase = models.PodPending
	pod.NodeName = ""

	if err := s.store.CreatePod(&pod); err != nil {
		log.Printf("Error creating pod %s/%s: %v", pod.Namespace, pod.Name, err)
		if errors.Is(err, store.ErrPodExists) {
			c.JSON(409, gin.H{"error": "Failed to create pod", "detail": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to create pod", "detail": err.Error()})
		}
		return
	}
	log.Printf("Created pod %s/%s successfully", pod.Namespace, pod.Name)

	s.watchManager.Publish(namespace, models.WatchEvent{
		EventType:   "ADDED",
		EventObject: "pod",
		Pod:         &pod,
	})

	c.JSON(201, pod)
}

func (s *APIServer) getPodHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("podname")
	pod, err := s.store.GetPod(namespace, name)
	if err != nil {
		c.JSON(404, gin.H{"error": "Pod not found", "detail": err.Error()})
	}
	c.JSON(200, pod)
}

func (s *APIServer) updatePodHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	originalName := c.Param("podname")

	var pod models.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "detail": err.Error()})
		return
	}

	if pod.Name == "" {
		c.JSON(400, gin.H{"error": "A pod name must be provided"})
		return
	}

	if pod.Namespace != namespace {
		c.JSON(400, gin.H{"error": fmt.Sprintf("No pod named %s in specified namespace %s", pod.Name, namespace)})
		return
	}

	if _, err := s.store.GetPod(namespace, originalName); err != nil {
		c.JSON(404, gin.H{"error": "Pod does not exist", "detail": err.Error()})
		return
	}

	if err := s.store.UpdatePod(&pod); err != nil {
		log.Printf("Failed to update pod: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update pod", "detail": err.Error()})
		return
	}
	log.Printf("Updated pod %s/%s successfully", pod.Namespace, pod.Name)

	s.watchManager.Publish(namespace, models.WatchEvent{
		EventType:   "MODIFIED",
		EventObject: "pod",
		Pod:         &pod,
	})

	c.JSON(200, pod)
}

func (s *APIServer) deletePodHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("podname")

	if err := s.store.DeletePod(namespace, name); err != nil {
		log.Printf("Error deleting pod %s/%s: %v", namespace, name, err)
		if errors.Is(err, store.ErrPodNotExist) {
			c.JSON(404, gin.H{"error": "Pod not found for deletion", "detail": err.Error()})
		} else {
			c.JSON(409, gin.H{"error": "Pod is already being deleted", "detail": err.Error()})
		}
		return
	}
	log.Printf("Pod %s/%s successfuly set for deletion", namespace, name)

	s.watchManager.Publish(namespace, models.WatchEvent{
		EventType:   "DELETED",
		EventObject: "pod",
		Pod: &models.Pod{
			Name:      name,
			Namespace: namespace,
		},
	})

	c.JSON(200, gin.H{"message": fmt.Sprintf("Pod %s/%s successfully set for deletion", namespace, name)})
}

func (s *APIServer) listPodsHandler(c *gin.Context) {
	watch := c.Query("watch")

	if watch == "true" {
		s.watchPods(c)
		return
	}

	namespace := c.Param("namespace")
	podList, err := s.store.ListPods(namespace)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not fetch pod list", "detail": err.Error()})
		return
	}

	c.JSON(200, podList)
}

func (s *APIServer) watchPods(c *gin.Context) {
	namespace := c.Param("namespace")

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.WriteHeader(200)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed assertion"})
	}

	watchCh := s.watchManager.Subscribe(namespace)
	defer s.watchManager.Unsubscribe(namespace, watchCh)

	ctx := c.Request.Context()

	for {
		select {
		case event := <-watchCh:
			if err := json.NewEncoder(c.Writer).Encode(event); err != nil {
				log.Printf("Error encoding watch event: %v", err)
				return
			}
			flusher.Flush()

		case <-ctx.Done():
			log.Printf("Client connection closed for namespace %s", namespace)
			return
		}
	}
}
