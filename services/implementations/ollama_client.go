package implementations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type OllamaClient struct {
	baseURL string
	model   string
	client  *http.Client
}

type OllamaRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  OllamaOptions   `json:"options,omitempty"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	MaxTokens   int     `json:"num_predict,omitempty"` // Ollama usa num_predict
}

type OllamaResponse struct {
	Model           string        `json:"model"`
	CreatedAt       string        `json:"created_at"`
	Message         OllamaMessage `json:"message"`
	Done            bool          `json:"done"`
	TotalDuration   int64         `json:"total_duration,omitempty"`
	PromptEvalCount int           `json:"prompt_eval_count,omitempty"`
	EvalCount       int           `json:"eval_count,omitempty"`
}

func NewOllamaClient() (*OllamaClient, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("OLLAMA_BASE_URL environment variable is required")
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		return nil, fmt.Errorf("OLLAMA_MODEL environment variable is required")
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second, // Ollama pode ser mais lento para gerar código
		},
	}, nil
}

func NewOllamaClientWithConfig(baseURL, model string) (*OllamaClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}
	if model == "" {
		return nil, fmt.Errorf("model is required")
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (c *OllamaClient) GenerateCode(ctx context.Context, prompt, contextInfo string) (string, error) {
	log.Printf("[INFO] Starting Ollama code generation with model: %s", c.model)

	systemPrompt := `You are an expert code generator. Generate production-ready code following best practices.
Return ONLY valid JSON with file paths as keys and file contents as values.
Example: {"main.go": "package main\\n\\nfunc main() {\\n\\tprintln(\\"Hello\\")\\n}", "README.md": "# Project\\n\\nDescription"}
Do not include any other text or explanations outside the JSON.`

	fullPrompt := fmt.Sprintf("Project Context: %s\n\nUser Request: %s", contextInfo, prompt)

	request := OllamaRequest{
		Model: c.model,
		Messages: []OllamaMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fullPrompt},
		},
		Stream: false,
		Options: OllamaOptions{
			Temperature: 0.7, // Balanço entre criatividade e consistência
			TopP:        0.9,
			MaxTokens:   4096,
		},
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("[DEBUG] Sending request to Ollama: %s", c.baseURL+"/api/chat")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	if ollamaResp.Message.Content == "" {
		return "", fmt.Errorf("no response from Ollama")
	}

	log.Printf("[INFO] Ollama generation completed. Tokens used - Prompt: %d, Generated: %d",
		ollamaResp.PromptEvalCount, ollamaResp.EvalCount)

	return ollamaResp.Message.Content, nil
}

func (c *OllamaClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check returned status %d", resp.StatusCode)
	}

	log.Printf("[INFO] Ollama is healthy at %s", c.baseURL)
	return nil
}

func (c *OllamaClient) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list models request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode models list: %w", err)
	}

	models := make([]string, len(result.Models))
	for i, model := range result.Models {
		models[i] = model.Name
	}

	return models, nil
}
