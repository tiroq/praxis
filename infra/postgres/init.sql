-- Praxis scaffold database initialization.
-- TODO: replace with migrations once the domain model stabilizes.

CREATE TABLE IF NOT EXISTS scaffold_marker (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO scaffold_marker (name)
VALUES ('praxis_initialized')
ON CONFLICT (name) DO NOTHING;
