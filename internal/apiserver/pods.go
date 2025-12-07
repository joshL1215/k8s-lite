package apiserver

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joshL1215/k8s-lite/internal/models"
	"github.com/joshL1215/k8s-lite/internal/store"
)

func (s *APIServer) createPodHandler(c *gin.Context) {
	var pod models.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "detail": err.Error()})
		return
	}

	if pod.Name == "" {
		c.JSON(400, gin.H{"error": "A pod name must be provided"})
		return
	}
	namespace := c.Param("namespace")
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
}
