package main

import (
	"github.com/joshL1215/k8s-lite/internal/apiserver"
	"github.com/joshL1215/k8s-lite/internal/store/memory"
)

const DefaultPort = "8080"

func main() {
	dataStore := memory.CreateInMemoryStore()
	apiServer := apiserver.CreateAPIServer(dataStore)
	apiServer.Serve(":" + DefaultPort)
}
