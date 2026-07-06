package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client encapsulates the concrete HTTP protocol to the llm-router service.
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// GenerateReplyInput carries semantic caller input for reply generation.
type GenerateReplyInput struct {
	InputEventID  string                `json:"input_event_id,omitempty"`
	CorrelationID string                `json:"correlation_id,omitempty"`
	Source        string                `json:"source,omitempty"`
	UserMessage   string                `json:"user_message"`
	Conversation  []ConversationMessage `json:"conversation,omitempty"`
	Metadata      map[string]string     `json:"metadata,omitempty"`
}

// ConversationMessage is an optional prior turn passed to the router.
type ConversationMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type generateReplyResponse struct {
	ReplyText      string `json:"reply_text"`
	AssistantReply string `json:"assistant_reply"`
}

// NewClient constructs a concrete llm-router client.
func NewClient(endpoint string, httpClient *http.Client) *Client {
	client := httpClient
	if client == nil {
		client = &http.Client{}
	}

	return &Client{
		endpoint:   strings.TrimSpace(endpoint),
		httpClient: client,
	}
}

// GenerateReply requests a reply from the llm-router service.
func (c *Client) GenerateReply(ctx context.Context, input GenerateReplyInput) (string, error) {
	if c == nil || c.endpoint == "" {
		return "", fmt.Errorf("llm router endpoint is not configured")
	}

	input.UserMessage = strings.TrimSpace(input.UserMessage)
	if input.UserMessage == "" {
		return "", fmt.Errorf("user_message is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("marshal llm reply request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build llm router request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call llm router: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("llm router returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var reply generateReplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return "", fmt.Errorf("decode llm router response: %w", err)
	}

	text := strings.TrimSpace(reply.ReplyText)
	if text == "" {
		text = strings.TrimSpace(reply.AssistantReply)
	}
	if text == "" {
		return "", fmt.Errorf("llm router response missing reply_text")
	}

	return text, nil
}
