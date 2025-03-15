package logic

import (
	"context"
	"fmt"
	"log"
	"sync"

	"node/clients"
	"node/models"
)

// GetBestHost finds the best available host for a given model
func GetBestHost(modelName string, redis *clients.RedisClient, logger *log.Logger) (*models.LLMHost, error) {
	logger.Printf("Fetching best host for model: %s", modelName)

	// Get all models from Redis
	allModels, err := redis.GetAllLLModels(context.Background())
	if err != nil {
		logger.Printf("Error fetching models from Redis: %v", err)
		return nil, err
	}

	// Find the requested model
	var selectedModel *models.LLModel
	for _, model := range allModels {
		if model.Modelinfo.Name == modelName {
			selectedModel = &model
			break
		}
	}

	if selectedModel == nil {
		logger.Printf("Model not found: %s", modelName)
		return nil, fmt.Errorf("model not found")
	}

	// Filter active hosts
	activeHosts := []string{}
	for _, hostingServer := range selectedModel.HostingServers {
		if hostingServer.Status {
			activeHosts = append(activeHosts, hostingServer.HostName)
		}
	}

	logger.Printf("Active hosts found: %v", activeHosts)

	// Case 1: If there's only one active host, return it immediately
	if len(activeHosts) == 1 {
		host, err := redis.GetLLMHost(context.Background(), activeHosts[0])
		if err != nil {
			logger.Printf("Error fetching host details for %s: %v", activeHosts[0], err)
			return nil, err
		}
		logger.Printf("Returning single active host: %s", host.HostName)
		return host, nil
	}

	// Case 2: Multiple active hosts, find the one with the lowest task count
	if len(activeHosts) > 1 {
		logger.Println("Fetching details of multiple active hosts...")

		var wg sync.WaitGroup
		hostChan := make(chan *models.LLMHost, len(activeHosts))
		errChan := make(chan error, len(activeHosts))

		for _, ip := range activeHosts {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				host, err := redis.GetLLMHost(context.Background(), ip)
				if err != nil {
					logger.Printf("Failed to fetch host %s: %v", ip, err)
					errChan <- err
					return
				}
				hostChan <- host
			}(ip)
		}

		wg.Wait()
		close(hostChan)
		close(errChan)

		// Pick the host with the lowest task count
		var bestHost *models.LLMHost
		minTaskCount := 5 // Max int value

		for host := range hostChan {
			logger.Printf("Host %s has task count: %d", host.HostName, len(host.Tasks))
			if len(host.Tasks) < minTaskCount {
				bestHost = host
				minTaskCount = len(host.Tasks)
			}
		}

		if bestHost != nil {
			logger.Printf("Selected best host: %s with task count: %d", bestHost.HostName, len(bestHost.Tasks))
			return bestHost, nil
		}
	}

	logger.Println("No available hosts found")
	return nil, fmt.Errorf("no available hosts found")
}

// GetServerToLoad finds an inactive hosting server for a given model
func GetServerToLoad(modelName string, redis *clients.RedisClient, logger *log.Logger) (*models.LLMHost, *models.LLModel, error) {
	logger.Printf("Finding inactive server for model: %s", modelName)

	// Get all models from Redis
	allModels, err := redis.GetAllLLModels(context.Background())
	if err != nil {
		logger.Printf("Error fetching models from Redis: %v", err)
		return nil, nil, err
	}

	// Find the requested model
	var selectedModel *models.LLModel
	for _, model := range allModels {
		if model.Modelinfo.Name == modelName {
			selectedModel = &model
			break
		}
	}

	if selectedModel == nil {
		logger.Printf("Model not found: %s", modelName)
		return nil, nil, fmt.Errorf("model not found")
	}

	// Find the first inactive server
	for _, hostingServer := range selectedModel.HostingServers {
		if !hostingServer.Status { // Inactive server found
			logger.Printf("Inactive server found: %s", hostingServer.HostName)

			// Retrieve server details from Redis
			host, err := redis.GetLLMHost(context.Background(), hostingServer.HostName)
			if err != nil {
				logger.Printf("Error fetching inactive host details (%s): %v", hostingServer.HostName, err)
				return nil, nil, err
			}

			logger.Printf("Returning inactive host: %s", host.HostName)
			return host, selectedModel, nil
		}
	}

	logger.Println("No inactive servers found")
	return nil, nil, fmt.Errorf("no inactive servers available")
}
