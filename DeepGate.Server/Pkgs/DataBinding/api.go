package databinding

const NodePort = "8080"
const HostPort = "9090"

// Message represents a single message in a chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletion represents a chat completion request
type ChatCompletion struct {
	Model    string                 `json:"model"`
	Messages []Message              `json:"messages"`
	Format   map[string]interface{} `json:"format"`
	TaskId   string                 `json:"task_id"`
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
