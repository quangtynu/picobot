package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// OllamaProvider talks to a local Ollama/vLLM OpenAI-compatible endpoint (no API key required).
type OllamaProvider struct {
	APIBase string
	Client  *http.Client
}

func NewOllamaProvider(apiBase string) *OllamaProvider {
	if apiBase == "" {
		apiBase = "http://localhost:11434/v1"
	}
	return &OllamaProvider{APIBase: strings.TrimRight(apiBase, "/"), Client: &http.Client{Timeout: 60 * time.Second}}
}

func (p *OllamaProvider) GetDefaultModel() string { return "meta-llama/Llama-3.1-8B-Instruct" }

// Reuse same minimal request/response shapes as OpenRouter provider (OpenAI-compatible)

type ollamaChatRequest struct {
	Model    string        `json:"model"`
	Messages []messageJSON `json:"messages"`
	Tools    []toolWrapper `json:"tools,omitempty"`
}

// Chat implements LLMProvider.Chat
func (p *OllamaProvider) Chat(ctx context.Context, messages []Message, tools []ToolDefinition, model string) (LLMResponse, error) {
	if model == "" {
		model = p.GetDefaultModel()
	}
	reqBody := ollamaChatRequest{Model: model, Messages: make([]messageJSON, 0, len(messages))}
	for _, m := range messages {
		mj := messageJSON{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallID}
		for _, tc := range m.ToolCalls {
			argsBytes, _ := json.Marshal(tc.Arguments)
			mj.ToolCalls = append(mj.ToolCalls, toolCallJSON{
				ID:   tc.ID,
				Type: "function",
				Function: toolCallFunctionJSON{
					Name:      tc.Name,
					Arguments: string(argsBytes),
				},
			})
		}
		reqBody.Messages = append(reqBody.Messages, mj)
	}

	// Include tools in modern format if provided
	if len(tools) > 0 {
		reqBody.Tools = make([]toolWrapper, 0, len(tools))
		for _, t := range tools {
			params := t.Parameters
			if params == nil {
				params = map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}
			}
			reqBody.Tools = append(reqBody.Tools, toolWrapper{
				Type: "function",
				Function: functionDef{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  params,
				},
			})
		}
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return LLMResponse{}, err
	}
	url := fmt.Sprintf("%s/chat/completions", p.APIBase)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(b)))
	if err != nil {
		return LLMResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Ollama typically does not require Authorization header when talking to localhost

	resp, err := p.Client.Do(req)
	if err != nil {
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// attempt to read response body for more details (do not expose credentials)
		bodyBytes, _ := io.ReadAll(resp.Body)
		body := strings.TrimSpace(string(bodyBytes))
		log.Printf("Ollama API non-2xx: %s body=%q", resp.Status, body)
		if body == "" {
			return LLMResponse{}, fmt.Errorf("Ollama API error: %s", resp.Status)
		}
		return LLMResponse{}, fmt.Errorf("Ollama API error: %s - %s", resp.Status, body)
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return LLMResponse{}, err
	}

	if len(out.Choices) == 0 {
		return LLMResponse{}, errors.New("Ollama API returned no choices")
	}

	msg := out.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		var tcs []ToolCall
		for _, tc := range msg.ToolCalls {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &parsed); err != nil {
				continue
			}
			tcs = append(tcs, ToolCall{ID: tc.ID, Name: tc.Function.Name, Arguments: parsed})
		}
		if len(tcs) > 0 {
			return LLMResponse{Content: strings.TrimSpace(msg.Content), HasToolCalls: true, ToolCalls: tcs}, nil
		}
	}

	return LLMResponse{Content: strings.TrimSpace(msg.Content), HasToolCalls: false}, nil
}
