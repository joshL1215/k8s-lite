package apiserver

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joshL1215/k8s-lite/pkg/store"
)

type APIServer struct {
	router *gin.Engine
	store  store.StoreInterface // having an interface here makes it store-implementation-agnostic
}

func (s *APIServer) Serve(port string) {
	log.Println("Serving API server...")
	if err := s.router.Run(port); err != nil {
		log.Printf("Could not serve API server: %v", err)
	}
}

// func (s *APIServer) registerRoutes() {
// 	podsGroup := s.router.Group("/api/v1/namespace/:namespace/pods") // version APIs for backwards compatability
// 	{

// 	}
// }

func CreateAPIServer(s store.StoreInterface) *APIServer {
	apiServer := &APIServer{
		router: gin.Default(),
		store:  s,
	}
	// apiServer.registerRoutes()
	return apiServer
}
