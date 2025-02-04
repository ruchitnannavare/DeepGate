package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	databinding "Pkgs/DataBinding"

	"github.com/gin-gonic/gin"
)

type NodeServer struct {
	logger             *log.Logger
	databaseConnection *databinding.DatabaseConnections
	hostIP             string
}

func NewNodeServer() *NodeServer {
	logger := databinding.ConfigureLogger()

	dbConnections, err := databinding.InitializeDatabases()
	if err != nil {
		logger.Fatalf("Failed to initialize databases: %v", err)
	}

	return &NodeServer{
		logger:             logger,
		databaseConnection: dbConnections,
	}
}

func (ns *NodeServer) SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		var infoPackage databinding.InfoPackage
		if err := c.BindJSON(&infoPackage); err != nil {
			ns.logger.Printf("Invalid ping request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Store in MongoDB
		collection := ns.databaseConnection.MongoDB.Collection("pings")
		_, err := collection.InsertOne(context.Background(), infoPackage)
		if err != nil {
			ns.logger.Printf("Failed to store ping in MongoDB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage failed"})
			return
		}

		// Cache in Redis
		jsonData, _ := json.Marshal(infoPackage)
		err = ns.databaseConnection.RedisClient.Set(
			context.Background(),
			fmt.Sprintf("ping:%s", infoPackage.IPAddress),
			jsonData,
			24*time.Hour,
		).Err()
		if err != nil {
			ns.logger.Printf("Failed to cache ping in Redis: %v", err)
		}

		ns.logger.Printf("Received ping from %s", infoPackage.IPAddress)
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

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
