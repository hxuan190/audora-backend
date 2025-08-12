-- +goose Up
-- +goose StatementBegin

-- Upload sessions table for tracking file uploads
CREATE TABLE upload_sessions (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    artist_id BIGINT NOT NULL, -- No FK reference
    user_id BIGINT NOT NULL, -- No FK reference
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    object_path TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for upload_sessions table
CREATE INDEX idx_upload_sessions_artist_id ON upload_sessions(artist_id);
CREATE INDEX idx_upload_sessions_user_id ON upload_sessions(user_id);
CREATE INDEX idx_upload_sessions_status ON upload_sessions(status);
CREATE INDEX idx_upload_sessions_expires_at ON upload_sessions(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_upload_sessions_expires_at;
DROP INDEX IF EXISTS idx_upload_sessions_status;
DROP INDEX IF EXISTS idx_upload_sessions_user_id;
DROP INDEX IF EXISTS idx_upload_sessions_artist_id;
DROP TABLE IF EXISTS upload_sessions;

-- +goose StatementEnd
