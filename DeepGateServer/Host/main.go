package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"./discovery"
	"github.com/gin-gonic/gin"
)

type HostStatus struct {
	HostID string `json:"host_id"`
	Status string `json:"status"`
}

func main() {
	// Initialize service discovery
	sd := discovery.NewServiceDiscovery("load-balancer", 8080)
	if err := sd.Advertise(context.Background()); err != nil {
		log.Fatalf("Failed to start service discovery: %v", err)
	}

	r := gin.Default()

	// Get Host B's address from environment variables or configuration
	hostBAddress := os.Getenv("HOST_B_ADDRESS")
	if hostBAddress == "" {
		log.Fatal("HOST_B_ADDRESS environment variable is not set")
	}

	// Unique identifier for Host C
	hostID := "host_c_1" // This should be unique for each instance of Host C

	// Register with Host B
	registerWithHostB(hostBAddress, hostID)

	// Wait for further requests (this is a placeholder for actual request handling logic)
	select {}
}

func registerWithHostB(hostBAddress, hostID string) {
	hostStatus := HostStatus{
		HostID: hostID,
		Status: "alive",
	}

	data, err := json.Marshal(hostStatus)
	if err != nil {
		log.Fatalf("Failed to marshal host status: %v", err)
	}

	url := fmt.Sprintf("http://%s/ping", hostBAddress)
	for {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Failed to register with Host B: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Println("Successfully registered with Host B")
				break
			} else {
				log.Printf("Failed to register with Host B, status code: %d", resp.StatusCode)
			}
		}

		// Retry after a delay
		time.Sleep(10 * time.Second)
	}
}
