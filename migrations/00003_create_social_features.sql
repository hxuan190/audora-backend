-- +goose Up
-- +goose StatementBegin

-- Tips/payments
CREATE TABLE tips (
    id UUID PRIMARY KEY NOT NULL,
    from_user_id UUID NOT NULL, -- No FK reference
    to_artist_id UUID NOT NULL, -- No FK reference
    song_id UUID, -- No FK reference, optional: tip for specific song
    amount_cents INTEGER NOT NULL, -- Store in cents to avoid decimal issues
    currency VARCHAR(3) DEFAULT 'USD',
    stripe_payment_intent_id VARCHAR(100) UNIQUE,
    stripe_charge_id VARCHAR(100),
    platform_fee_cents INTEGER NOT NULL,
    artist_payout_cents INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, completed, failed, refunded
    message TEXT, -- Optional message from tipper
    is_anonymous BOOLEAN DEFAULT false,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Real-time artist-fan messaging
CREATE TABLE artist_messages (
    id UUID PRIMARY KEY NOT NULL,
    artist_id UUID NOT NULL, -- No FK reference
    message_text TEXT NOT NULL,
    target_type VARCHAR(50) NOT NULL, -- all_active_listeners, specific_song_listeners, followers
    target_song_id UUID, -- No FK reference
    sent_to_count INTEGER DEFAULT 0,
    read_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Message deliveries (tracking who received messages)
CREATE TABLE message_deliveries (
    id UUID PRIMARY KEY NOT NULL,
    message_id UUID NOT NULL, -- No FK reference
    user_id UUID NOT NULL, -- No FK reference
    delivered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(message_id, user_id)
);

-- Current listening sessions (for real-time dashboard)
CREATE TABLE listening_sessions (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID, -- No FK reference
    session_id VARCHAR(100) NOT NULL, -- For anonymous users
    song_id UUID NOT NULL, -- No FK reference
    artist_id UUID NOT NULL, -- No FK reference
    country_code VARCHAR(2),
    city VARCHAR(100),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_heartbeat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Create indexes for listening_sessions table
CREATE INDEX idx_listening_sessions_artist_id ON listening_sessions(artist_id);
CREATE INDEX idx_listening_sessions_active ON listening_sessions(is_active, last_heartbeat);

-- User preferences
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL, -- No FK reference
    preferred_genres INTEGER[], -- Array of genre IDs
    preferred_moods INTEGER[], -- Array of mood IDs
    auto_play BOOLEAN DEFAULT true,
    shuffle_by_default BOOLEAN DEFAULT false,
    notification_new_releases BOOLEAN DEFAULT true,
    notification_artist_messages BOOLEAN DEFAULT true,
    notification_tips_received BOOLEAN DEFAULT true,
    explicit_content_allowed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_listening_sessions_active;
DROP INDEX IF EXISTS idx_listening_sessions_artist_id;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS listening_sessions;
DROP TABLE IF EXISTS message_deliveries;
DROP TABLE IF EXISTS artist_messages;
DROP TABLE IF EXISTS tips;

-- +goose StatementEnd 