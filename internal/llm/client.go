package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tiroq/praxis/internal/storage/userfacts"
)

// Client encapsulates the concrete HTTP protocol to the llm-router service.
type Client struct {
	endpoint   string
	httpClient *http.Client
}

type ConversationMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// GenerateReplyRequest carries semantic caller input for reply generation.
type GenerateReplyRequest struct {
	InputEventID        string              `json:"input_event_id,omitempty"`
	CorrelationID       string              `json:"correlation_id,omitempty"`
	Source              string              `json:"source,omitempty"`
	UserMessage         string              `json:"user_message"`
	ConversationContext ConversationContext `json:"-"`
	Metadata            map[string]string   `json:"metadata,omitempty"`
}

// ExtractFactsRequest carries semantic input for user fact extraction.
type ExtractFactsRequest struct {
	CorrelationID       string              `json:"correlation_id,omitempty"`
	Source              string              `json:"source,omitempty"`
	LatestUserMessage   string              `json:"latest_user_message"`
	ConversationContext ConversationContext `json:"-"`
	Metadata            map[string]string   `json:"metadata,omitempty"`
}

// FactExtractionResponse is the structured extraction output from llm-router.
type FactExtractionResponse struct {
	Facts []userfacts.CandidateFact
}

type generateReplyResponse struct {
	ReplyText      string `json:"reply_text"`
	AssistantReply string `json:"assistant_reply"`
}

type generateReplyPayload struct {
	InputEventID  string                `json:"input_event_id,omitempty"`
	CorrelationID string                `json:"correlation_id,omitempty"`
	Source        string                `json:"source,omitempty"`
	UserMessage   string                `json:"user_message"`
	Conversation  []ConversationMessage `json:"conversation,omitempty"`
	Metadata      map[string]string     `json:"metadata,omitempty"`
}

type extractFactsPayload struct {
	CorrelationID     string                `json:"correlation_id,omitempty"`
	Source            string                `json:"source,omitempty"`
	LatestUserMessage string                `json:"latest_user_message"`
	Conversation      []ConversationMessage `json:"conversation,omitempty"`
	Metadata          map[string]string     `json:"metadata,omitempty"`
}

type extractFactsResponse struct {
	Facts []extractedFact `json:"facts"`
}

type extractedFact struct {
	Type       string  `json:"type"`
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
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

// GenerateReply requests a semantic reply from the llm-router service.
func (c *Client) GenerateReply(ctx context.Context, input GenerateReplyRequest) (string, error) {
	if c == nil || c.endpoint == "" {
		return "", fmt.Errorf("llm router endpoint is not configured")
	}

	input.UserMessage = strings.TrimSpace(input.UserMessage)
	if input.UserMessage == "" {
		return "", fmt.Errorf("user_message is required")
	}

	payload, err := json.Marshal(generateReplyPayload{
		InputEventID:  input.InputEventID,
		CorrelationID: input.CorrelationID,
		Source:        input.Source,
		UserMessage:   input.UserMessage,
		Conversation:  contextMessages(input.ConversationContext),
		Metadata:      input.Metadata,
	})
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

// ExtractFacts requests structured candidate user facts from the llm-router service.
func (c *Client) ExtractFacts(ctx context.Context, input ExtractFactsRequest) (FactExtractionResponse, error) {
	if c == nil || c.endpoint == "" {
		return FactExtractionResponse{}, fmt.Errorf("llm router endpoint is not configured")
	}

	input.LatestUserMessage = strings.TrimSpace(input.LatestUserMessage)
	if input.LatestUserMessage == "" {
		return FactExtractionResponse{}, fmt.Errorf("latest_user_message is required")
	}

	payload, err := json.Marshal(extractFactsPayload{
		CorrelationID:     input.CorrelationID,
		Source:            input.Source,
		LatestUserMessage: input.LatestUserMessage,
		Conversation:      contextMessages(input.ConversationContext),
		Metadata:          input.Metadata,
	})
	if err != nil {
		return FactExtractionResponse{}, fmt.Errorf("marshal fact extraction request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.extractFactsEndpoint(), bytes.NewReader(payload))
	if err != nil {
		return FactExtractionResponse{}, fmt.Errorf("build fact extraction request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return FactExtractionResponse{}, fmt.Errorf("call fact extraction endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return FactExtractionResponse{}, fmt.Errorf("fact extraction returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var output extractFactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return FactExtractionResponse{}, fmt.Errorf("decode fact extraction response: %w", err)
	}

	facts := make([]userfacts.CandidateFact, 0, len(output.Facts))
	for _, fact := range output.Facts {
		factType := strings.TrimSpace(fact.Type)
		value := strings.TrimSpace(fact.Value)
		if factType == "" || value == "" {
			continue
		}
		if fact.Confidence < 0 || fact.Confidence > 1 {
			return FactExtractionResponse{}, fmt.Errorf("fact extraction response confidence outside [0,1]")
		}
		facts = append(facts, userfacts.CandidateFact{
			Type:       factType,
			Value:      value,
			Confidence: fact.Confidence,
		})
	}

	return FactExtractionResponse{Facts: facts}, nil
}

func contextMessages(ctx ConversationContext) []ConversationMessage {
	if len(ctx.Messages) == 0 {
		return nil
	}
	return ctx.Messages
}

func (c *Client) extractFactsEndpoint() string {
	parsed, err := url.Parse(c.endpoint)
	if err != nil {
		return strings.TrimRight(c.endpoint, "/") + "/extract-facts"
	}
	if strings.HasSuffix(parsed.Path, "/v1/reply") {
		parsed.Path = strings.TrimSuffix(parsed.Path, "/v1/reply") + "/v1/extract-facts"
		return parsed.String()
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/extract-facts"
	return parsed.String()
}
