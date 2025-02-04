package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	databinding "Pkgs/DataBinding"

	"github.com/gin-gonic/gin"
)

type HostServer struct {
	logger   *log.Logger
	nodeIP   string
	hostName string
}

func NewHostServer() *HostServer {
	hostName, _ := os.Hostname()
	return &HostServer{
		logger:   databinding.ConfigureLogger(),
		hostName: hostName,
	}
}

func (hs *HostServer) ScanNetwork() {
	hs.logger.Println("Starting network scan...")

	// Get local network IP range
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		hs.logger.Printf("Error getting network interfaces: %v", err)
		return
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				baseIP := strings.Join(strings.Split(ipnet.IP.String(), ".")[:3], ".") + "."

				for i := 1; i <= 254; i++ {
					testIP := fmt.Sprintf("%s%d", baseIP, i)

					// Skip the current machine's IP
					if testIP == ipnet.IP.String() {
						continue
					}

					go hs.tryPingNode(testIP)
				}
			}
		}
	}
}

func (hs *HostServer) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (hs *HostServer) tryPingNode(ip string) {
	url := fmt.Sprintf("http://%s:8080/ping", ip)

	infoPackage := databinding.InfoPackage{
		IPAddress:  hs.getLocalIP(),
		Identifier: 0, // Host
		HostName:   hs.hostName,
		Timestamp:  time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(infoPackage)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		hs.logger.Printf("Ping to %s failed: %v", ip, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		hs.logger.Printf("Node found at %s!", ip)
		hs.nodeIP = ip
		// Optionally, you might want to stop further scanning here
		return
	}
}

func (hs *HostServer) SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Add any local HTTP routes as needed
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"node_ip": hs.nodeIP,
			"status":  "ready",
		})
	})

	return r
}

func (hs *HostServer) Run() {
	// First scan the network
	hs.ScanNetwork()

	r := hs.SetupRoutes()
	hs.logger.Println("Host server starting on localhost:9090")
	r.Run("localhost:9090")
}

func main() {
	hostServer := NewHostServer()
	hostServer.Run()
}
