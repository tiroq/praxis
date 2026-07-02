package sqlite

const schema = `
CREATE TABLE IF NOT EXISTS events (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	source TEXT NOT NULL,
	subject_id TEXT NOT NULL,
	correlation_id TEXT,
	causation_id TEXT,
	trace_id TEXT,
	occurred_at TEXT NOT NULL,
	payload TEXT NOT NULL,
	metadata TEXT NOT NULL,
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);
CREATE INDEX IF NOT EXISTS idx_events_subject_id ON events(subject_id);
CREATE INDEX IF NOT EXISTS idx_events_correlation_id ON events(correlation_id);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at ON events(occurred_at);
`
