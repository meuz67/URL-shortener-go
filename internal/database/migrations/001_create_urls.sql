CREATE TABLE IF NOT EXISTS urls (
    code TEXT PRIMARY KEY,
    encrypted_url BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS urls_created_at_idx ON urls(created_at);
