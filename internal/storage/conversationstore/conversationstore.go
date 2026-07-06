package conversationstore

// ListFilter specifies criteria for filtering messages in a conversation.
type ListFilter struct {
	Limit  int // maximum number of messages to return (0 = default)
	Offset int // number of messages to skip
}
