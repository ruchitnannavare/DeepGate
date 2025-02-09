package services

import (
	"encoding/json"

	"host/clients"
)

// APIService struct for interacting with multiple APIs
type APIService struct {
	LocalAPI *clients.APIClient
}

// NewAPIService initializes a new API service
func NewAPIService(localBaseURL string) *APIService {
	return &APIService{
		LocalAPI: clients.NewAPIClient(localBaseURL),
	}
}

// FetchLocalData calls the local API
func (s *APIService) FetchLocalModelList() (map[string]interface{}, error) {
	respBody, err := s.LocalAPI.MakeRequest("GET", "/api/tags", nil, nil)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
