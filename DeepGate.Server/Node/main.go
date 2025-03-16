package main

import (
	"log"
	"net"

	databinding "Pkgs/DataBinding"

	"node/clients"
	"node/routes"

	"github.com/gin-gonic/gin"
)

type NodeServer struct {
	logger             *log.Logger
	databaseConnection *databinding.DatabaseConnections
	APIRepo            *clients.APIClient
	redis              *clients.RedisClient
}

func NewNodeServer() *NodeServer {
	logger := databinding.ConfigureLogger()

	dbConnections, err := databinding.InitializeDatabases()
	if err != nil {
		logger.Fatalf("Failed to initialize databases: %v", err)
	}

	redis := clients.NewRedisClient(dbConnections.RedisClient)

	return &NodeServer{
		logger:             logger,
		databaseConnection: dbConnections,
		redis:              redis,
	}
}

func (ns *NodeServer) getLocalIP() string {
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

func (ns *NodeServer) SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Initialize and register host routes
	// Host routing logic
	hostHandler := routes.NewHostHandler(ns.logger, ns.redis)
	hostHandler.RegisterRoutes(r)

	// Client routing logic
	clientHandler := routes.NewClientHandler(ns.logger, ns.redis)
	clientHandler.RegisterRoutes(r)

	return r
}

func (ns *NodeServer) Run() {
	r := ns.SetupRoutes()
	nodeIp := ns.getLocalIP()
	ns.logger.Println("Find node at ip: " + nodeIp)
	ns.logger.Println("Node server starting on 0.0.0.0:8080")
	r.Run("0.0.0.0:8080")
}

func main() {
	nodeServer := NewNodeServer()
	nodeServer.Run()
}
