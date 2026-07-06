package natsworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tiroq/praxis/internal/core/kernel"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

// LLMRouterClient calls the dedicated LLM router service for assistant replies.
type LLMRouterClient struct {
	endpoint   string
	httpClient *http.Client
}

// NewLLMRouterClient creates a router client for the given endpoint.
func NewLLMRouterClient(endpoint string, httpClient *http.Client) *LLMRouterClient {
	client := httpClient
	if client == nil {
		client = &http.Client{}
	}
	return &LLMRouterClient{
		endpoint:   strings.TrimSpace(endpoint),
		httpClient: client,
	}
}

type llmReplyRequest struct {
	InputEventID   string                `json:"input_event_id"`
	CorrelationID  string                `json:"correlation_id"`
	Source         string                `json:"source"`
	InputText      string                `json:"input_text"`
	Metadata       map[string]string     `json:"metadata,omitempty"`
	PipelineResult kernel.PipelineResult `json:"pipeline_result"`
}

type llmReplyResponse struct {
	AssistantReply string `json:"assistant_reply"`
}

// GenerateReply requests an assistant reply from the router.
func (c *LLMRouterClient) GenerateReply(ctx context.Context, input natstransport.InputMessage, result kernel.PipelineResult) (string, error) {
	if c == nil || c.endpoint == "" {
		return "", fmt.Errorf("llm router endpoint is not configured")
	}

	reqBody := llmReplyRequest{
		InputEventID:   input.ID,
		CorrelationID:  input.CorrelationID,
		Source:         input.Source,
		InputText:      input.Text,
		Metadata:       input.Metadata,
		PipelineResult: result,
	}

	payload, err := json.Marshal(reqBody)
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

	var reply llmReplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return "", fmt.Errorf("decode llm router response: %w", err)
	}

	reply.AssistantReply = strings.TrimSpace(reply.AssistantReply)
	if reply.AssistantReply == "" {
		return "", fmt.Errorf("llm router response missing assistant_reply")
	}

	return reply.AssistantReply, nil
}
