package apiserver

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joshL1215/k8s-lite/internal/store"
)

const DefaultNamespace = "default"

type APIServer struct {
	router       *gin.Engine
	store        store.StoreInterface // having an interface here makes it store-implementation-agnostic
	watchManager watchManager
}

func (s *APIServer) Serve(port string) {
	log.Println("Serving API server...")
	if err := s.router.Run(port); err != nil {
		log.Printf("Could not serve API server: %v", err)
	}
}

func (s *APIServer) registerRoutes() {
	podsGroup := s.router.Group("/api/v1/namespace/:namespace/pods") // version APIs for backwards compatability
	{
		podsGroup.POST("", s.createPodHandler)
		podsGroup.GET("", s.listPodsHandler) // includes a query parameter ?watch= to open a long lived TCP connection for watching
		podsGroup.GET("/:podname", s.getPodHandler)
		podsGroup.PUT("/:podname", s.updatePodHandler)
		podsGroup.DELETE(":podname", s.deletePodHandler)
	}

	nodesGroup := s.router.Group("/api/v1/nodes")
	{
		nodesGroup.POST("", s.createNodeHandler)
		nodesGroup.GET("", s.listNodesHandler)
		nodesGroup.GET("/:nodename", s.getNodeHandler)
		nodesGroup.PUT("/:nodename", s.updateNodeHandler)
		nodesGroup.DELETE("/:nodename", s.deleteNodeHandler)
	}
}

func CreateAPIServer(s store.StoreInterface) *APIServer {
	apiServer := &APIServer{
		router:       gin.Default(),
		store:        s,
		watchManager: *NewWatchManager(),
	}
	apiServer.registerRoutes()
	return apiServer
}
