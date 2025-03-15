package main

import (
	"log"

	databinding "Pkgs/DataBinding"

	"node/clients"
	"node/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type NodeServer struct {
	logger             *log.Logger
	databaseConnection *databinding.DatabaseConnections
	hostIP             string
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

func (ns *NodeServer) SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Initialize and register host routes
	// Host routing logic
	hostHandler := routes.NewHostHandler(ns.logger, ns.redis)
	hostHandler.RegisterRoutes(r)

	// Client routing logic
	clientHandler := routes.NewClientHandler(ns.logger, ns.redis)
	clientHandler.RegisterRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (ns *NodeServer) Run() {
	r := ns.SetupRoutes()
	ns.logger.Println("Node server starting on 0.0.0.0:8080")
	r.Run("0.0.0.0:8080")
}

func main() {
	nodeServer := NewNodeServer()
	nodeServer.Run()
}
