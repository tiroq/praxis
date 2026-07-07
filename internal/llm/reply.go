package llm

import "context"

// ReplyService owns reply-specific orchestration details.
type ReplyService struct {
	client *Client
}

func NewReplyService(client *Client) *ReplyService {
	return &ReplyService{client: client}
}

type ReplyRequest struct {
	InputEventID        string
	CorrelationID       string
	Source              string
	UserMessage         string
	Metadata            map[string]string
	ConversationContext ConversationContext
}

func (s *ReplyService) Generate(ctx context.Context, req ReplyRequest) (string, error) {
	return s.client.GenerateReply(ctx, GenerateReplyRequest{
		InputEventID:        req.InputEventID,
		CorrelationID:       req.CorrelationID,
		Source:              req.Source,
		UserMessage:         req.UserMessage,
		Metadata:            req.Metadata,
		ConversationContext: req.ConversationContext,
	})
}
