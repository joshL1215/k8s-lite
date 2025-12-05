package main

import (
	"github.com/joshL1215/k8s-lite/cmd/apiserver"
	"github.com/joshL1215/k8s-lite/pkg/store/memory"
)

const DefaultPort = "8080"

func main() {
	dataStore := memory.CreateInMemoryStore()
	apiServer := apiserver.CreateAPIServer(dataStore)
	apiServer.Serve(":" + DefaultPort)
}
