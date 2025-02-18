package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"node/models"

	"github.com/go-redis/redis/v8"
)

const (
	LLModelsKey      = "llm_models" // Key for storing list of models
	LLMHostKeyPrefix = "llm_host:"  // Prefix for host keys
	DefaultTTL       = 24 * time.Hour
)

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient now accepts the existing Redis client from DatabaseConnections
func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
	}
}

// AddOrUpdateLLModel now takes a slice of models and saves them directly
func (rc *RedisClient) UpdateLLModelList(ctx context.Context, models []models.LLModel) error {
	data, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal models: %v", err)
	}

	return rc.client.Set(ctx, LLModelsKey, data, DefaultTTL).Err()
}

// 🔹 AddOrUpdateLLModel() - Adds a new model or updates an existing one
func (rc *RedisClient) AddOrUpdateLLModel(ctx context.Context, newModel models.LLModel) error {
	// Get existing models
	allModels, err := rc.GetAllLLModels(ctx)
	if err != nil {
		return err
	}

	// Check if model exists and update it
	found := false
	for i, model := range allModels {
		if model.Modelinfo.Name == newModel.Modelinfo.Name {
			allModels[i] = newModel // Update existing model
			found = true
			break
		}
	}

	// If not found, append it
	if !found {
		allModels = append(allModels, newModel)
	}

	// Save updated models back to Redis
	return rc.UpdateLLModelList(ctx, allModels)
}

// GetAllLLModels retrieves all LLModels from Redis
func (rc *RedisClient) GetAllLLModels(ctx context.Context) ([]models.LLModel, error) {
	data, err := rc.client.Get(ctx, LLModelsKey).Bytes()
	if err == redis.Nil {
		return []models.LLModel{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get models from Redis: %v", err)
	}

	var models []models.LLModel
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, fmt.Errorf("failed to unmarshal models: %v", err)
	}

	return models, nil
}

// RemoveLLModel removes a model from the Redis list
func (rc *RedisClient) RemoveLLModel(ctx context.Context, modelName string) error {
	all_models, err := rc.GetAllLLModels(ctx)
	if err != nil {
		return err
	}

	filtered := make([]models.LLModel, 0)
	for _, model := range all_models {
		if model.Modelinfo.Name != modelName {
			filtered = append(filtered, model)
		}
	}

	return rc.saveModels(ctx, filtered)
}

// saveModels helper function to save models to Redis
func (rc *RedisClient) saveModels(ctx context.Context, models []models.LLModel) error {
	data, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal models: %v", err)
	}

	return rc.client.Set(ctx, LLModelsKey, data, DefaultTTL).Err()
}

// SaveLLMHost stores a host with its IP as the key
func (rc *RedisClient) SaveLLMHost(ctx context.Context, host models.LLMHost) error {
	data, err := json.Marshal(host)
	if err != nil {
		return fmt.Errorf("failed to marshal host: %v", err)
	}

	key := LLMHostKeyPrefix + host.HostInfo.IPAddress
	return rc.client.Set(ctx, key, data, DefaultTTL).Err()
}

// GetLLMHost retrieves a host by its IP address
func (rc *RedisClient) GetLLMHost(ctx context.Context, ipAddress string) (*models.LLMHost, error) {
	key := LLMHostKeyPrefix + ipAddress
	data, err := rc.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get host from Redis: %v", err)
	}

	var host models.LLMHost
	if err := json.Unmarshal(data, &host); err != nil {
		return nil, fmt.Errorf("failed to unmarshal host: %v", err)
	}

	return &host, nil
}

// RemoveLLMHostByIP removes a host by its IP address
func (rc *RedisClient) RemoveLLMHostByIP(ctx context.Context, ipAddress string) error {
	key := LLMHostKeyPrefix + ipAddress
	return rc.client.Del(ctx, key).Err()
}

// GetAllLLMHostIPs retrieves all host IP addresses
func (rc *RedisClient) GetAllLLMHostIPs(ctx context.Context) ([]string, error) {
	pattern := LLMHostKeyPrefix + "*"
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get host keys: %v", err)
	}

	ips := make([]string, len(keys))
	for i, key := range keys {
		ips[i] = key[len(LLMHostKeyPrefix):] // Remove prefix to get IP
	}

	return ips, nil
}
