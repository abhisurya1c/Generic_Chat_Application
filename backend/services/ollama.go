package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const OllamaAPIURL = "http://localhost:11434/api/generate"

// GenerateRequest represents the body for Ollama API
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse represents the non-streaming response from Ollama
type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// StreamResponse represents the streaming response chunk from Ollama
type StreamResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// CallOllama handles both streaming and non-streaming requests to Ollama
func CallOllama(prompt string, model string, stream bool, onChunk func(string) error) (string, error) {
	reqBody := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: stream,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", OllamaAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 0, // No timeout for streaming
	}
	if !stream {
		client.Timeout = 120 * time.Second
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned status: %d", resp.StatusCode)
	}

	if !stream {
		var fullResp GenerateResponse
		if err := json.NewDecoder(resp.Body).Decode(&fullResp); err != nil {
			return "", fmt.Errorf("failed to decode response: %v", err)
		}
		return fullResp.Response, nil
	}

	// Handle Streaming
	scanner := bufio.NewScanner(resp.Body)
	var fullResponseBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			continue // Skip invalid lines
		}

		if onChunk != nil {
			if err := onChunk(streamResp.Response); err != nil {
				return fullResponseBuilder.String(), err
			}
		}

		fullResponseBuilder.WriteString(streamResp.Response)

		if streamResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fullResponseBuilder.String(), fmt.Errorf("streaming error: %v", err)
	}

	return fullResponseBuilder.String(), nil
}
