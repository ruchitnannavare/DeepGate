package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// APIClient struct to manage API calls
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *log.Logger
}

// NewAPIClient initializes a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, // Set timeout for API calls
		},
	}
}

func MakeTemporaryAPIClient(ipAddress, hostPort string) *APIClient {
	return NewAPIClient(fmt.Sprintf("http://%s:%s", ipAddress, hostPort))
}

// MakeRequest is a generic method to send API requests
func (c *APIClient) MakeRequest(method, endpoint string, body interface{}, headers map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		var ok = string(jsonBody)
		fmt.Println(ok)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create a new request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(respBody))
	}

	return respBody, nil
}

// In clients/api_client.go

func (c *APIClient) MakeStreamingRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + endpoint
	c.HTTPClient = &http.Client{
		Timeout: 0, // No timeout for streaming
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	return c.HTTPClient.Do(req)
}
