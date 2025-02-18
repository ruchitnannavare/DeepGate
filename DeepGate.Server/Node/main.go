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
	// r.POST("/ping", func(c *gin.Context) {
	// 	var infoPackage databinding.InfoPackage
	// 	if err := c.BindJSON(&infoPackage); err != nil {
	// 		ns.logger.Printf("Invalid ping request: %v", err)
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 		return
	// 	}

	// 	ns.APIRepo = clients.NewAPIClient(fmt.Sprintf("http://%s:%s", infoPackage.IPAddress, infoPackage.HostPort))

	// 	resp, err := ns.APIRepo.MakeRequest("GET", "/fetchlocalmodellist", nil, nil)
	// 	if err != nil {
	// 		return
	// 	}

	// 	// After parsing the original model response
	// 	modelResponse, err := models.ParseModelResponse(resp)
	// 	if err != nil {
	// 		ns.logger.Printf("Failed to parse model response: %v", err)
	// 		return
	// 	}

	// 	// Convert to simplified host model info
	// 	hostModels := models.ConvertModelsToHostInfo(modelResponse.Models)

	// 	// Create LLMHost object to maintain host and model information
	// 	llmHost := models.LLMHost{
	// 		HostInfo:  infoPackage,
	// 		ModelInfo: hostModels,
	// 	}

	// 	// Add or update the host in Redis
	// 	err = ns.redis.AddOrUpdateLLMHost(context.Background(), llmHost)
	// 	if err != nil {
	// 		ns.logger.Printf("Failed to update LLMHost in Redis: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update host information"})
	// 		return
	// 	}

	// 	ns.logger.Printf("Received ping from %s", infoPackage.IPAddress)
	// 	c.JSON(http.StatusOK, gin.H{"status": "received"})
	// })

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
