package databinding

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InfoPackage represents the handshake information for network discovery
type InfoPackage struct {
	IPAddress  string `json:"ip_address"`
	Identifier int    `json:"identifier"` // 0 for host, 1 for client
	HostName   string `json:"host_name"`
	Timestamp  int64  `json:"timestamp"`
	HostPort   string `json:"host_port"`
}

// DatabaseConnections holds connections for MongoDB and Redis
type DatabaseConnections struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
	MongoDB     *mongo.Database
}

// InitializeDatabases sets up MongoDB and Redis connections
func InitializeDatabases() (*DatabaseConnections, error) {
	// MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := "mongodb://localhost:27017"
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return nil, err
	}

	// Redis Connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test Redis connection
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		return nil, err
	}

	// Initialize database
	database := mongoClient.Database("network_discovery")

	return &DatabaseConnections{
		MongoClient: mongoClient,
		RedisClient: redisClient,
		MongoDB:     database,
	}, nil
}

// ConfigureLogger sets up a custom logger with timestamps
func ConfigureLogger() *log.Logger {
	return log.New(os.Stdout, "[DEEPGATE] ", log.Ldate|log.Ltime|log.Lshortfile)
}
