package databinding

// Message represents a single message in a chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletion represents a chat completion request
type ChatCompletion struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// StreamResponse represents the structure of each streaming response chunk
type StreamResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}
