package clients

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	databinding "Pkgs/DataBinding"
)

// APIService struct for interacting with multiple APIs
type OllamaClient struct {
	api    *APIClient
	logger *log.Logger
}

// NewAPIService initializes a new API service with logging
func NewOllamaClient(ollama_port string, logger *log.Logger) *OllamaClient {
	return &OllamaClient{
		api:    NewAPIClient(fmt.Sprintf("http://localhost:%s", ollama_port)),
		logger: logger,
	}
}

// FetchLocalModelList calls the local API to get available models
func (o *OllamaClient) FetchLocalModelList() (map[string]interface{}, error) {
	start := time.Now()
	o.logger.Printf("Fetching local model list...")

	respBody, err := o.api.MakeRequest("GET", "/api/tags", nil, nil)
	if err != nil {
		o.logger.Printf("Error fetching model list: %v", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		o.logger.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	elapsed := time.Since(start)
	o.logger.Printf("Successfully fetched model list in %v", elapsed)
	return data, nil
}

// LoadLocalModel loads a specific model
func (o *OllamaClient) LoadLocalModel(modelName string) error {
	start := time.Now()
	o.logger.Printf("Starting to load model: %s", modelName)

	jsonPayload := map[string]interface{}{
		"model": modelName,
	}

	respBody, err := o.api.MakeRequest("POST", "/api/generate", jsonPayload, nil)
	if err != nil {
		o.logger.Printf("Error loading model %s: %v", modelName, err)
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		o.logger.Printf("Error parsing response: %v", err)
		return err
	}

	elapsed := time.Since(start)
	o.logger.Printf("Successfully loaded model %s in %v", modelName, elapsed)
	return nil
}

// StopModel stops a running model instance
func (o *OllamaClient) StopModel() error {
	start := time.Now()
	o.logger.Printf("Attempting to stop running model...")

	// Ollama's stop endpoint doesn't require a model name as it stops the currently running model
	respBody, err := o.api.MakeRequest("POST", "/api/stop", nil, nil)
	if err != nil {
		o.logger.Printf("Error stopping model: %v", err)
		return fmt.Errorf("failed to stop model: %v", err)
	}

	// Check if we got a response
	if len(respBody) > 0 {
		var response map[string]interface{}
		if err := json.Unmarshal(respBody, &response); err != nil {
			o.logger.Printf("Warning: Could not parse stop response: %v", err)
			// Don't return error as the stop might have still succeeded
		}
	}

	elapsed := time.Since(start)
	o.logger.Printf("Successfully stopped model in %v", elapsed)
	return nil
}

// CheckModelStatus checks if a model is loaded and ready
func (o *OllamaClient) CheckModelStatus(modelName string) (bool, error) {
	start := time.Now()
	o.logger.Printf("Checking status for model: %s", modelName)

	models, err := o.FetchLocalModelList()
	if err != nil {
		o.logger.Printf("Error checking model status: %v", err)
		return false, err
	}

	// Check if model exists in the list
	modelList, ok := models["models"].([]interface{})
	if !ok {
		o.logger.Printf("Invalid response format while checking model status")
		return false, fmt.Errorf("invalid response format")
	}

	for _, model := range modelList {
		if modelMap, ok := model.(map[string]interface{}); ok {
			if modelMap["name"] == modelName {
				elapsed := time.Since(start)
				o.logger.Printf("Model %s status check completed in %v", modelName, elapsed)
				return true, nil
			}
		}
	}

	elapsed := time.Since(start)
	o.logger.Printf("Model %s not found (checked in %v)", modelName, elapsed)
	return false, nil
}

// StreamChatCompletion handles streaming chat completion requests
func (o *OllamaClient) StreamChatCompletion(chat databinding.ChatCompletion, responseChan chan string, errorChan chan error) {
	start := time.Now()
	o.logger.Printf("Starting chat completion streaming for model: %s", chat.Model)

	// Prepare the request payload
	jsonData, err := json.Marshal(chat)
	if err != nil {
		o.logger.Printf("Error marshaling chat completion request: %v", err)
		errorChan <- err
		return
	}

	// Make the streaming request
	resp, err := o.api.MakeStreamingRequest("POST", "/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		o.logger.Printf("Error initiating streaming request: %v", err)
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	// Create a scanner to read the streaming response
	scanner := bufio.NewScanner(resp.Body)

	// Process each chunk of the streaming response
	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println("Received line:", line)

		var streamResp databinding.StreamResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			o.logger.Printf("Error parsing stream chunk: %v", err)
			continue
		}

		// Send the content through the channel
		if streamResp.Message.Content != "" {
			responseChan <- streamResp.Message.Content
		}

		// If done, close the channel
		if streamResp.Done {
			elapsed := time.Since(start)
			o.logger.Printf("Streaming completed in %v", elapsed)
			close(responseChan)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		o.logger.Printf("Error reading stream: %v", err)
		errorChan <- err
	}
}
