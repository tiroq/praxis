package conversationstore

const schema = `
-- Conversations table: projection grouping derived from correlation metadata.
-- Rebuildable from immutable events; not canonical truth.
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    correlation_id TEXT NOT NULL UNIQUE,
    lifecycle TEXT NOT NULL CHECK(lifecycle IN ('created', 'active', 'archived')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_message_at TEXT
);

-- Create index for correlation_id lookups
CREATE INDEX IF NOT EXISTS idx_conversations_correlation_id ON conversations(correlation_id);

-- Messages table: append-only conversation history projection (RFC-033 Projection Store)
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('user', 'assistant')),
    content TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    metadata TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id)
);

-- Create indexes for message queries
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
CREATE INDEX IF NOT EXISTS idx_messages_event_id ON messages(event_id);
`
