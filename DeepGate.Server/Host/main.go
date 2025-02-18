package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	databinding "Pkgs/DataBinding"

	"host/clients"
	"host/routes"

	"github.com/gin-gonic/gin"
)

type HostServer struct {
	logger   *log.Logger
	nodeIP   string
	hostName string
	ollama   *clients.OllamaClient
}

func NewHostServer() *HostServer {
	hostName, _ := os.Hostname()
	host_logger := databinding.ConfigureLogger()
	ollamaClient := clients.NewOllamaClient("11434", host_logger)
	return &HostServer{
		logger:   host_logger,
		hostName: hostName,
		ollama:   ollamaClient,
	}
}

func (hs *HostServer) ScanNetwork() {
	hs.logger.Println("Starting network scan...")

	ips := hs.getActiveIPs()
	if len(ips) == 0 {
		hs.logger.Println("No active IP addresses found.")
		return
	}

	localIP := hs.getLocalIP()
	for _, ip := range ips {
		// Skip the current machine's IP
		if ip == localIP {
			continue
		}

		go hs.tryPingNode(ip)
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
		HostPort:   "9090",
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

func (hs *HostServer) getActiveIPs() []string {
	var ips []string

	switch runtime.GOOS {
	case "linux":
		ips = hs.getActiveIPsLinux()
	case "windows":
		ips = hs.getActiveIPsWindows()
	case "darwin":
		ips = hs.getActiveIPsDarwin()
	default:
		hs.logger.Printf("Unsupported OS: %s", runtime.GOOS)
	}

	return ips
}

func (hs *HostServer) getActiveIPsLinux() []string {
	var ips []string

	arpFile := "/proc/net/arp"
	file, err := os.Open(arpFile)
	if err != nil {
		hs.logger.Printf("Error opening %s: %v", arpFile, err)
		return ips
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isFirstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if isFirstLine {
			// Skip header line
			isFirstLine = false
			continue
		}

		// Split the line into fields
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		ip := fields[0]
		ips = append(ips, ip)
	}

	if err := scanner.Err(); err != nil {
		hs.logger.Printf("Error reading %s: %v", arpFile, err)
	}

	return ips
}

func (hs *HostServer) getActiveIPsWindows() []string {
	var ips []string

	cmd := exec.Command("arp", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		hs.logger.Printf("Error executing arp command: %v", err)
		return ips
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "Interface:") || strings.HasPrefix(line, "Internet Address") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 1 {
			ip := fields[0]
			// Validate that it's an IP address
			if net.ParseIP(ip) != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips
}

func (hs *HostServer) getActiveIPsDarwin() []string {
	var ips []string

	cmd := exec.Command("arp", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		hs.logger.Printf("Error executing arp command: %v", err)
		return ips
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Line format: ? (192.168.1.1) at 0:1a:2b:3c:4d:5e on en0 ifscope [ethernet]
		// Extract IP inside parentheses

		start := strings.Index(line, "(")
		end := strings.Index(line, ")")
		if start != -1 && end != -1 && start < end {
			ip := line[start+1 : end]
			if net.ParseIP(ip) != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips
}

func (hs *HostServer) SetupRoutes() *gin.Engine {
	r := gin.Default()

	routeHandler := routes.NewRouteHandler(hs.logger, hs.ollama)
	routeHandler.RegisterRoutes(r)

	return r
}

func (hs *HostServer) Run() {
	// First scan the network
	hs.ScanNetwork()

	r := hs.SetupRoutes()
	hs.logger.Println("Host server starting on 0.0.0.0:9090")
	r.Run("0.0.0.0:9090")
}

func main() {
	hostServer := NewHostServer()
	hostServer.Run()
}
