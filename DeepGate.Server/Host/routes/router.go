package routes

import (
	"host/clients"
	"io"
	"log"
	"net/http"
	"time"

	databinding "Pkgs/DataBinding"

	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	logger *log.Logger
	ollama *clients.OllamaClient
}

func NewRouteHandler(logger *log.Logger, ollamaClient *clients.OllamaClient) *RouteHandler {
	return &RouteHandler{
		logger: logger,
		ollama: ollamaClient,
	}
}

// RegisterRoutes registers all host-related routes
func (r *RouteHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/host/load-model", r.handleLoadModel)
	router.GET("/host/fetch-models", r.handleFetchLocalModelList)
	router.POST("/host/chat", r.handleChatCompletion)
}

func (r *RouteHandler) handleLoadModel(c *gin.Context) {
	var request struct {
		ModelName string `json:"model_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	start := time.Now()
	r.logger.Printf("Received request to load model: %s", request.ModelName)

	err := r.ollama.LoadLocalModel(request.ModelName)
	if err != nil {
		r.logger.Printf("Failed to load model %s: %v", request.ModelName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	elapsed := time.Since(start)
	r.logger.Printf("Successfully loaded model %s in %v", request.ModelName, elapsed)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Model loaded successfully",
		"model":      request.ModelName,
		"time_taken": elapsed.String(),
	})
}

func (r *RouteHandler) handleFetchLocalModelList(c *gin.Context) {
	start := time.Now()
	r.logger.Printf("Fetching local model list")

	models, err := r.ollama.FetchLocalModelList()
	if err != nil {
		r.logger.Printf("Error fetching model list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	elapsed := time.Since(start)
	r.logger.Printf("Successfully fetched model list in %v", elapsed)

	c.JSON(http.StatusOK, models)
}

func (r *RouteHandler) handleChatCompletion(c *gin.Context) {
	var chatRequest databinding.ChatCompletion

	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		r.logger.Printf("Error binding chat request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid chat request format",
		})
		return
	}

	start := time.Now()
	r.logger.Printf("Starting chat completion with model: %s", chatRequest.Model)

	// Create channels for streaming
	responseChan := make(chan string)
	errorChan := make(chan error)

	// Set up Server-Sent Events
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// Start streaming in a goroutine
	go r.ollama.StreamChatCompletion(chatRequest, responseChan, errorChan)

	// Stream the response
	c.Stream(func(w io.Writer) bool {
		select {
		case content, ok := <-responseChan:
			if !ok {
				elapsed := time.Since(start)
				r.logger.Printf("Chat completion finished in %v", elapsed)
				return false
			}
			c.SSEvent("message", content)
			return true
		case err := <-errorChan:
			r.logger.Printf("Error during chat completion: %v", err)
			c.SSEvent("error", err.Error())
			return false
		case <-c.Request.Context().Done():
			r.logger.Printf("Client disconnected")
			return false
		}
	})
}
