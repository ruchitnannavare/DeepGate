package routes

import (
	databinding "Pkgs/DataBinding"
	"context"
	"log"
	"net/http"
	"node/clients"
	_ "node/docs"
	"node/models"

	"github.com/gin-gonic/gin"
)

type HostHandler struct {
	logger *log.Logger
	redis  *clients.RedisClient
}

func NewHostHandler(logger *log.Logger, redis *clients.RedisClient) *HostHandler {
	return &HostHandler{
		logger: logger,
		redis:  redis,
	}
}

// RegisterRoutes registers all host-related routes
func (h *HostHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/ping", h.handlePing)
	router.POST("/node/complete-task", h.handleCompleteTask)
}

// handleCompleteTask handles the complete task request from hosts
func (h *HostHandler) handleCompleteTask(c *gin.Context) {
	var request struct {
		TaskId   string `json:"task_id" binding:"required"`
		HostName string `json:"host_name" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		h.logger.Printf("Invalid task request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	host, err := h.redis.GetLLMHost(context.Background(), request.HostName)
	if err != nil {
		h.logger.Printf("Failed to get LLMHost from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get host information"})
		return
	}

	// Filter out the TaskId from host.Tasks
	filteredTasks := make([]string, 0, len(host.Tasks))
	for _, task := range host.Tasks {
		if task != request.TaskId {
			filteredTasks = append(filteredTasks, task)
		}
	}
	host.Tasks = filteredTasks

	h.redis.SaveLLMHost(context.Background(), *host)

	c.JSON(http.StatusOK, gin.H{"status": "task " + request.TaskId + " completed"})
}

// handlePing handles the ping request from hosts
func (h *HostHandler) handlePing(c *gin.Context) {
	var infoPackage databinding.InfoPackage
	if err := c.BindJSON(&infoPackage); err != nil {
		h.logger.Printf("Invalid ping request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Make temporary API client to make a request
	apiClient := clients.MakeTemporaryAPIClient(infoPackage.IPAddress, infoPackage.HostPort)

	resp, err := apiClient.MakeRequest("GET", "/fetchlocalmodellist", nil, nil)
	if err != nil {
		h.logger.Printf("Failed to fetch model list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch model list"})
		return
	}

	// Parse model response
	modelResponse, err := models.ParseModelResponse(resp)
	if err != nil {
		h.logger.Printf("Failed to parse model response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse model response"})
		return
	}

	// Convert to simplified host model info
	hostModels := models.ConvertModelsToHostInfo(modelResponse.Models)

	// Create LLMHost object to maintain host and model information
	llmHost := models.LLMHost{
		HostName:  infoPackage.HostName,
		HostInfo:  infoPackage,
		ModelInfo: hostModels,
	}

	// Add or update the host in Redis
	err = h.redis.SaveLLMHost(context.Background(), llmHost)
	if err != nil {
		h.logger.Printf("Failed to update LLMHost in Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update host information"})
		return
	}

	// Get current list of LLModels
	allModels, err := h.redis.GetAllLLModels(context.Background())
	if err != nil {
		h.logger.Printf("Failed to get LLModels from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model information"})
		return
	}

	// Create a map for faster lookup of existing models
	modelMap := make(map[string]*models.LLModel)
	for i, model := range allModels {
		modelMap[model.Modelinfo.Name] = &allModels[i]
	}

	// Process each model from the host
	for _, hostModel := range hostModels {
		if existingModel, exists := modelMap[hostModel.Name]; exists {
			// Update existing model's host list
			hostExists := false
			for _, host := range existingModel.HostingServers {
				if host.HostName == infoPackage.HostName {
					hostExists = true
					break
				}
			}
			if !hostExists {
				hostingServer := models.HostingServer{
					HostName: infoPackage.HostName,
					Status:   false,
				}
				existingModel.HostingServers = append(existingModel.HostingServers, hostingServer)
			}
		} else {
			hostingServer := models.HostingServer{
				HostName: infoPackage.HostName,
				Status:   false,
			}
			newModel := models.LLModel{
				Modelinfo:      hostModel,
				HostingServers: []models.HostingServer{hostingServer},
			}
			allModels = append(allModels, newModel)
		}
	}

	// Save updated LLModels to Redis
	if err := h.redis.UpdateLLModelList(context.Background(), allModels); err != nil {
		h.logger.Printf("Failed to save updated LLModels to Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model information"})
		return
	}

	h.logger.Printf("Received ping from %s", infoPackage.HostName)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
