// Package routes defines API endpoints for the client and host services.
//
//	@title DeepGate API
//	@version 1.0
//	@description This is the API documentation for the DeepGate AI service.
//	@BasePath /
//
//	@contact.name API Support
//	@contact.email support@deepgate.com
package routes

import (
	"context"
	"io"
	"log"
	"net/http"
	"node/clients"
	_ "node/docs"
	"node/logic"

	databinding "Pkgs/DataBinding"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	logger *log.Logger
	redis  *clients.RedisClient
}

func NewClientHandler(logger *log.Logger, redis *clients.RedisClient) *ClientHandler {
	return &ClientHandler{
		logger: logger,
		redis:  redis,
	}
}

// RegisterRoutes registers all client-related routes
func (c *ClientHandler) RegisterRoutes(router *gin.Engine) {
	// @Summary Fetch available models
	// @Description Retrieves a list of all available AI models from Redis
	// @Tags Client
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /client/fetch-models [get]
	router.GET("/client/fetch-models", c.handleFetchModels)

	// @Summary Load a model
	// @Description Loads an AI model on an available host server
	// @Tags Client
	// @Accept json
	// @Produce json
	// @Param request body struct{Model string} true "Model to load"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /client/load-model [post]
	router.POST("/client/load-model", c.handleClientLoadModel)

	// @Summary Chat with AI model
	// @Description Sends a chat request to an AI model hosted on the best available server
	// @Tags Client
	// @Accept json
	// @Produce json
	// @Param request body databinding.ChatCompletion true "Chat request payload"
	// @Success 200 {object} string "Streaming response"
	// @Failure 400 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /client/chat [post]
	router.POST("/client/chat", c.handleClientChat)
}

// handleFetchModels fetches available models from Redis
func (c *ClientHandler) handleFetchModels(gc *gin.Context) {
	ctx := gc.Request.Context()

	// Fetch models from Redis
	models, err := c.redis.GetAllLLModels(ctx)
	if err != nil {
		c.logger.Println("Error fetching models from Redis:", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch models"})
		return
	}

	// Return models as JSON response
	gc.JSON(http.StatusOK, gin.H{"models": models})
}

// handleClientLoadModel loads a model on an available host
func (c *ClientHandler) handleClientLoadModel(gc *gin.Context) {
	var request struct {
		Model string `json:"model" binding:"required"`
	}

	if err := gc.ShouldBindJSON(&request); err != nil {
		c.logger.Printf("Invalid request format: %v", err)
		gc.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get inactive host and corresponding model
	inactiveHost, selectedModel, err := logic.GetServerToLoad(request.Model, c.redis, c.logger)
	if err != nil {
		c.logger.Printf("Failed to find inactive host: %v", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "No available host"})
		return
	}

	// Call the Host server to load model
	apiClient := clients.MakeTemporaryAPIClient(inactiveHost.HostInfo.IPAddress, inactiveHost.HostInfo.HostPort)

	resp, err := apiClient.MakeRequest("POST", "/load-model", request, nil)
	if err != nil {
		c.logger.Printf("Failed to load model: %v", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load model"})
		return
	}

	// Update HostingServer status to active
	for i, server := range selectedModel.HostingServers {
		if server.IPAdd == inactiveHost.HostInfo.IPAddress {
			selectedModel.HostingServers[i].Status = true
			c.logger.Printf("Updated server %s status to active", server.IPAdd)
			break
		}
	}

	// Save updated model back to Redis
	err = c.redis.AddOrUpdateLLModel(context.Background(), *selectedModel)
	if err != nil {
		c.logger.Printf("Failed to update model in Redis: %v", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model status"})
		return
	}

	c.logger.Printf("Model %s successfully loaded on host %s", request.Model, inactiveHost.HostInfo.IPAddress)

	// Forward the response
	gc.JSON(http.StatusOK, resp)
}

// handleClientChat handles chat requests with AI models
func (c *ClientHandler) handleClientChat(gc *gin.Context) {
	var chatRequest databinding.ChatCompletion

	if err := gc.ShouldBindJSON(&chatRequest); err != nil {
		c.logger.Printf("Invalid chat request format: %v", err)
		gc.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat request format"})
		return
	}

	// Get best available host
	bestHost, err := logic.GetBestHost(chatRequest.Model, c.redis, c.logger)
	if err != nil {
		c.logger.Printf("Failed to find best host: %v", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "No available host"})
		return
	}

	// Call Host server
	apiClient := clients.MakeTemporaryAPIClient(bestHost.HostInfo.IPAddress, bestHost.HostInfo.HostPort)

	// Open a streaming connection to the Host's /chat API
	hostResp, err := apiClient.MakeStreamRequest("POST", "/chat", nil, chatRequest)
	if err != nil {
		c.logger.Printf("Chat request failed: %v", err)
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Chat request failed"})
		return
	}
	defer hostResp.Body.Close()

	// Set headers for SSE
	gc.Header("Content-Type", "text/event-stream")
	gc.Header("Cache-Control", "no-cache")
	gc.Header("Connection", "keep-alive")
	gc.Header("Transfer-Encoding", "chunked")

	// Stream the response from Host to Client
	gc.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		for {
			n, err := hostResp.Body.Read(buffer)
			if n > 0 {
				gc.Writer.Write(buffer[:n])
				gc.Writer.Flush()
			}
			if err != nil {
				if err == io.EOF {
					return false // End of stream
				}
				c.logger.Printf("Error reading stream: %v", err)
				return false
			}
		}
	})
}
