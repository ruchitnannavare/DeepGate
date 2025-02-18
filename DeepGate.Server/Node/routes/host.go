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
	// @Summary Host Ping
	// @Description Receives a ping request from a host and updates model information
	// @Tags hosts
	// @Accept  json
	// @Produce  json
	// @Param infoPackage body databinding.InfoPackage true "Host information package"
	// @Success 200 {object} map[string]string "status: received"
	// @Failure 400 {object} map[string]string "error: Invalid request"
	// @Failure 500 {object} map[string]string "error: Failed to fetch model list"
	// @Router /ping [post]
	router.POST("/ping", h.handlePing)
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
		IPAdd:     infoPackage.IPAddress,
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
				if host.IPAdd == infoPackage.IPAddress {
					hostExists = true
					break
				}
			}
			if !hostExists {
				hostingServer := models.HostingServer{
					IPAdd:  infoPackage.IPAddress,
					Status: false,
				}
				existingModel.HostingServers = append(existingModel.HostingServers, hostingServer)
			}
		} else {
			hostingServer := models.HostingServer{
				IPAdd:  infoPackage.IPAddress,
				Status: false,
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

	h.logger.Printf("Received ping from %s", infoPackage.IPAddress)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
