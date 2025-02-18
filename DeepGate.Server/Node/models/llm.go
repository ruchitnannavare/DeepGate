package models

import (
	"encoding/json"
	"fmt"

	databinding "Pkgs/DataBinding"
)

// ModelResponse represents the top-level JSON response for fetchlocalmodels api
// START
type HostModelInfoResponse struct {
	Models []Model `json:"models"`
}

// Model represents a single model's information
type Model struct {
	Details    ModelDetails `json:"details"`
	Digest     string       `json:"digest"`
	Model      string       `json:"model"`
	ModifiedAt string       `json:"modified_at"`
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
}

// ModelDetails contains the detailed information about a model
type ModelDetails struct {
	Families          []string `json:"families"`
	Family            string   `json:"family"`
	Format            string   `json:"format"`
	ParameterSize     string   `json:"parameter_size"`
	ParentModel       string   `json:"parent_model"`
	QuantizationLevel string   `json:"quantization_level"`
}

func ParseModelResponse(data []byte) (*HostModelInfoResponse, error) {
	var response HostModelInfoResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model response: %v", err)
	}
	return &response, nil
}

// END

// Base unit of deepgate-service-cluster
type LLMHost struct {
	IPAdd     string                  `json:"ip_add"`
	HostInfo  databinding.InfoPackage `json:"host_info"`
	ModelInfo []HostModelInfo         `json:"model_info"`
	Status    bool                    `json:"status"`
	TaskCount int                     `json:"task_count"`
}

type LLModel struct {
	Modelinfo      HostModelInfo
	HostingServers []HostingServer
}

type HostingServer struct {
	IPAdd  string
	Status bool
}

// HostModelInfo represents simplified model information with status
type HostModelInfo struct {
	Name          string `json:"name"`
	ParameterSize string `json:"parameter_size"`
	Family        string `json:"family"`
	Size          int64  `json:"size"`
}

// ConvertToHostModelInfo converts a Model to HostModelInfo
func ConvertToHostModelInfo(model Model) HostModelInfo {
	return HostModelInfo{
		Name:          model.Name,
		ParameterSize: model.Details.ParameterSize,
		Family:        model.Details.Family,
		Size:          model.Size, // Default task count
	}
}

// ConvertModelsToHostInfo converts a slice of Models to slice of HostModelInfo
func ConvertModelsToHostInfo(models []Model) []HostModelInfo {
	hostModels := make([]HostModelInfo, len(models))
	for i, model := range models {
		hostModels[i] = ConvertToHostModelInfo(model)
	}
	return hostModels
}
