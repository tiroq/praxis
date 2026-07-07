package userfacts

const schema = `
CREATE TABLE IF NOT EXISTS candidate_user_facts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    correlation_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    confidence REAL NOT NULL CHECK(confidence >= 0 AND confidence <= 1),
    source_event_id TEXT NOT NULL,
    source_message_id TEXT NOT NULL,
    validation_state TEXT NOT NULL CHECK(validation_state IN ('observed', 'extracted', 'correlated', 'reviewed', 'human_approved', 'canonical')),
    validation_updated_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS candidate_user_fact_validation_transitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fact_id TEXT NOT NULL,
    from_state TEXT NOT NULL,
    to_state TEXT NOT NULL CHECK(to_state IN ('observed', 'extracted', 'correlated', 'reviewed', 'human_approved', 'canonical')),
    actor TEXT NOT NULL,
    reason TEXT NOT NULL,
    transitioned_at TEXT NOT NULL,
    FOREIGN KEY (fact_id) REFERENCES candidate_user_facts(id)
);

CREATE INDEX IF NOT EXISTS idx_candidate_user_facts_user_id ON candidate_user_facts(user_id);
CREATE INDEX IF NOT EXISTS idx_candidate_user_facts_correlation_id ON candidate_user_facts(correlation_id);
CREATE INDEX IF NOT EXISTS idx_candidate_user_facts_source_event_id ON candidate_user_facts(source_event_id);
CREATE INDEX IF NOT EXISTS idx_candidate_user_facts_source_message_id ON candidate_user_facts(source_message_id);
CREATE INDEX IF NOT EXISTS idx_candidate_user_fact_validation_transitions_fact_id ON candidate_user_fact_validation_transitions(fact_id);
`
