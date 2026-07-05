CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ip_address INET NOT NULL,
    user_id INT,
    username VARCHAR(64),
    session_data JSONB NOT NULL DEFAULT '{}'::jsonb
);
