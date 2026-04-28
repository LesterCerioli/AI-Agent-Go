package implementations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DeepSeekClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type DeepSeekRequest struct {
	Model     string            `json:"model"`
	Messages  []DeepSeekMessage `json:"messages"`
	Stream    bool              `json:"stream"`
	MaxTokens int               `json:"max_tokens,omitempty"`
}

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	ID      string           `json:"id"`
	Choices []DeepSeekChoice `json:"choices"`
	Usage   DeepSeekUsage    `json:"usage"`
}

type DeepSeekChoice struct {
	Index        int             `json:"index"`
	Message      DeepSeekMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

type DeepSeekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewDeepSeekClient(apiKey, baseURL string) (*DeepSeekClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("DeepSeek API key is required")
	}

	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (c *DeepSeekClient) GenerateCode(ctx context.Context, prompt, contextInfo string) (string, error) {
	log.Printf("[INFO] Starting DeepSeek code generation")

	systemPrompt := `You are an expert code generator. Generate production-ready code following best practices.
Return ONLY valid JSON with file paths as keys and file contents as values.
Example: {"main.go": "package main\\n\\nfunc main() {\\n\\tprintln(\\"Hello\\")\\n}", "README.md": "# Project\\n\\nDescription"}
Do not include any other text or explanations outside the JSON.`

	fullPrompt := fmt.Sprintf("Project Context: %s\n\nUser Request: %s", contextInfo, prompt)

	request := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fullPrompt},
		},
		Stream:    false,
		MaxTokens: 4096,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var deepSeekResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&deepSeekResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", fmt.Errorf("no response from DeepSeek")
	}

	content := deepSeekResp.Choices[0].Message.Content
	log.Printf("[INFO] DeepSeek generation completed. Tokens used: %d", deepSeekResp.Usage.TotalTokens)

	return content, nil
}
