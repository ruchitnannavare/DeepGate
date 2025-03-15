package models

import (
	"log"

	"host/clients"
)

type HostServer struct {
	logger   *log.Logger
	nodeIP   string
	hostName string
	ollama   *clients.OllamaClient
}
