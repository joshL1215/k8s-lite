package apiserver

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joshL1215/k8s-lite/internal/api/models"
	"github.com/joshL1215/k8s-lite/internal/store"
)

func (s *APIServer) createNodeHandler(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "detail": err.Error()})
		return
	}

	if node.Name == "" {
		c.JSON(400, gin.H{"error": "A node name must be provided"})
		return
	}

	if node.Status == "" {
		node.Status = models.NodeReady
	}

	if err := s.store.CreateNode(&node); err != nil {
		log.Printf("Error creating node %s: %v", node.Name, err)
		if errors.Is(err, store.ErrNodeExists) {
			c.JSON(409, gin.H{"error": "Failed to create node", "detail": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to create node", "detail": err.Error()})
		}
		return
	}
	log.Printf("Created node %s successfully", node.Name)
	c.JSON(201, node)
}

func (s *APIServer) getNodeHandler(c *gin.Context) {
	name := c.Param("nodename")
	node, err := s.store.GetNode(name)
	if err != nil {
		c.JSON(404, gin.H{"error": "Node not found", "detail": err})
	}
	c.JSON(200, node)
}

func (s *APIServer) updateNodeHandler(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "detail": err.Error()})
		return
	}

	if node.Name == "" {
		c.JSON(400, gin.H{"error": "A node name must be provided"})
		return
	}

	if _, err := s.store.GetNode(node.Name); err != nil {
		if errors.Is(err, store.ErrNodeNotExist) {
			c.JSON(404, gin.H{"error": "Node does not exist", "detail": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to find pod", "detail": err.Error()})
		}
	}

	if err := s.store.UpdateNode(&node); err != nil {
		log.Printf("Failed to update node: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update node", "detail": err.Error()})
		return
	}

	c.JSON(200, node)
}

func (s *APIServer) deleteNodeHandler(c *gin.Context) {
	name := c.Param("nodename")

	if err := s.store.DeleteNode(name); err != nil {
		log.Printf("Error deleting node %s: %v", name, err)
		if errors.Is(err, store.ErrNodeNotExist) {
			c.JSON(404, gin.H{"error": "Node not found for deletion", "detail": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Unable to delete node", "detail": err.Error()})
		}
		return
	}

	log.Printf("Node %s successfully deleted", name)
	c.JSON(200, gin.H{"message": fmt.Sprintf("Node %s successfully deleted", name)})
}

func (s *APIServer) listNodesHandler(c *gin.Context) {
	nodeList, err := s.store.ListNodes()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to list nodes", "detail": err.Error()})
	}

	c.JSON(200, nodeList)
}
