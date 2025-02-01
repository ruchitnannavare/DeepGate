// host_b/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"GoPacks/discovery"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx             = context.Background()
	redisClient     *redis.Client
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
)

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust as needed
	})

	// Initialize MongoDB client
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	mongoCollection = mongoClient.Database("server_management").Collection("hosts")
}

type HostStatus struct {
	HostID string `json:"host_id" bson:"host_id"`
	Status string `json:"status" bson:"status"`
}

func main() {
	// Initialize service discovery
	sd := discovery.NewServiceDiscovery("load-balancer", 8080)
	if err := sd.Advertise(context.Background()); err != nil {
		log.Fatalf("Failed to start service discovery: %v", err)
	}

	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		var hostStatus HostStatus
		if err := c.ShouldBindJSON(&hostStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Store in Redis with a TTL
		err := redisClient.Set(ctx, hostStatus.HostID, hostStatus.Status, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store in Redis"})
			return
		}

		// Store in MongoDB
		_, err = mongoCollection.UpdateOne(
			ctx,
			bson.M{"host_id": hostStatus.HostID},
			bson.M{"$set": hostStatus},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store in MongoDB"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Host status updated"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
