// Package ai is the single backend client for OpenAI chat completions
// (docs/foundation-spec.md §6). The API key never leaves the backend.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.openai.com/v1"

// ErrNotConfigured means OPENAI_API_KEY is empty.
var ErrNotConfigured = errors.New("openai api key is not configured")

// Message is one chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Client calls the chat completions API. Safe for concurrent use.
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
	retryDelay time.Duration
}

// New creates a client. baseURL may be empty (OpenAI default); it exists so
// a compatible proxy or a local fake server can be used.
func New(apiKey string, model string, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		apiKey:     apiKey,
		model:      model,
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		retryDelay: 500 * time.Millisecond,
	}
}

// Configured reports whether an API key is set.
func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

// Schema describes a strict JSON-schema response format.
type Schema struct {
	Name   string
	Schema map[string]any
}

type responseFormat struct {
	Type       string            `json:"type"`
	JSONSchema *jsonSchemaFormat `json:"json_schema,omitempty"`
}

type jsonSchemaFormat struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
	Temperature    *float64        `json:"temperature,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// CreateChatCompletion returns the assistant's text reply.
func (c *Client) CreateChatCompletion(ctx context.Context, messages []Message) (string, error) {
	return c.complete(ctx, chatRequest{Model: c.model, Messages: messages})
}

// CompleteJSON asks for a reply matching schema and decodes it into out.
func (c *Client) CompleteJSON(ctx context.Context, messages []Message, schema Schema, out any) error {
	temperature := 0.0
	content, err := c.complete(ctx, chatRequest{
		Model:    c.model,
		Messages: messages,
		ResponseFormat: &responseFormat{
			Type:       "json_schema",
			JSONSchema: &jsonSchemaFormat{Name: schema.Name, Strict: true, Schema: schema.Schema},
		},
		Temperature: &temperature,
	})
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(content), out); err != nil {
		return fmt.Errorf("decode openai json reply: %w", err)
	}
	return nil
}

func (c *Client) complete(ctx context.Context, request chatRequest) (string, error) {
	if !c.Configured() {
		return "", ErrNotConfigured
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshal openai request: %w", err)
	}

	// One retry on 5xx or network errors.
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}
		content, retry, err := c.send(ctx, payload)
		if err == nil {
			return content, nil
		}
		lastErr = err
		if !retry {
			break
		}
	}
	return "", lastErr
}

func (c *Client) send(ctx context.Context, payload []byte) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", false, fmt.Errorf("build openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", ctx.Err() == nil, fmt.Errorf("call openai: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", true, fmt.Errorf("read openai response: %w", err)
	}
	retry := resp.StatusCode >= 500

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", retry, fmt.Errorf("parse openai response (status %d): %w", resp.StatusCode, err)
	}
	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return "", retry, fmt.Errorf("openai request failed (status %d): %s", resp.StatusCode, parsed.Error.Message)
		}
		return "", retry, fmt.Errorf("openai request failed with status %d", resp.StatusCode)
	}
	if len(parsed.Choices) == 0 {
		return "", false, errors.New("openai returned no choices")
	}
	return parsed.Choices[0].Message.Content, false, nil
}
