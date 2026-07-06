CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id INT,
    username VARCHAR(64),
    session_data JSONB NOT NULL DEFAULT '{}'::jsonb
);
