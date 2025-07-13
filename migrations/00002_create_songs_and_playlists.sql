-- +goose Up
-- +goose StatementBegin

-- Songs table
CREATE TABLE songs (
    id UUID PRIMARY KEY NOT NULL,
    artist_id UUID NOT NULL, -- No FK reference
    title VARCHAR(200) NOT NULL,
    description TEXT,
    file_url TEXT NOT NULL, -- S3/MinIO URL
    file_size_bytes BIGINT,
    duration_seconds INTEGER,
    artwork_url TEXT,
    genre_id UUID, -- No FK reference
    mood_id UUID, -- No FK reference
    tier content_tier NOT NULL DEFAULT 'public_discovery',
    ai_suggested_tier content_tier, -- AI recommendation
    tier_override_by_artist BOOLEAN DEFAULT false,
    bpm INTEGER,
    key_signature VARCHAR(10),
    is_explicit BOOLEAN DEFAULT false,
    is_processed BOOLEAN DEFAULT false,
    processing_status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    processing_error TEXT,
    play_count BIGINT DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    tip_count INTEGER DEFAULT 0,
    total_tips DECIMAL(10,2) DEFAULT 0.00,
    release_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Song plays tracking
CREATE TABLE song_plays (
    id UUID PRIMARY KEY NOT NULL,
    song_id UUID NOT NULL, -- No FK reference
    user_id UUID, -- No FK reference, nullable for anonymous plays
    session_id VARCHAR(100), -- For anonymous tracking
    ip_address INET,
    user_agent TEXT,
    country_code VARCHAR(2),
    city VARCHAR(100),
    duration_played_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT false, -- True if >80% played
    skip_reason VARCHAR(50), -- user_skip, next_song, etc.
    played_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for song_plays table
CREATE INDEX idx_song_plays_song_id ON song_plays(song_id);
CREATE INDEX idx_song_plays_user_id ON song_plays(user_id);
CREATE INDEX idx_song_plays_played_at ON song_plays(played_at);

-- User favorites/likes
CREATE TABLE user_favorites (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL, -- No FK reference
    song_id UUID NOT NULL, -- No FK reference
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, song_id)
);

-- Artist followers
CREATE TABLE artist_followers (
    id UUID PRIMARY KEY NOT NULL,
    artist_id UUID NOT NULL, -- No FK reference
    follower_user_id UUID NOT NULL, -- No FK reference
    notification_enabled BOOLEAN DEFAULT true,
    followed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(artist_id, follower_user_id)
);

-- Playlists
CREATE TABLE playlists (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    artwork_url TEXT,
    playlist_type VARCHAR(50) NOT NULL, -- curated, user_created, mood_based
    mood_id UUID, -- No FK reference
    created_by_user_id UUID, -- No FK reference
    is_public BOOLEAN DEFAULT true,
    play_count BIGINT DEFAULT 0,
    song_count INTEGER DEFAULT 0,
    total_duration_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Playlist songs (many-to-many)
CREATE TABLE playlist_songs (
    id UUID PRIMARY KEY NOT NULL,
    playlist_id UUID NOT NULL, -- No FK reference
    song_id UUID NOT NULL, -- No FK reference
    position INTEGER NOT NULL,
    added_by_user_id UUID, -- No FK reference
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(playlist_id, song_id),
    UNIQUE(playlist_id, position)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_song_plays_played_at;
DROP INDEX IF EXISTS idx_song_plays_user_id;
DROP INDEX IF EXISTS idx_song_plays_song_id;
DROP TABLE IF EXISTS playlist_songs;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS artist_followers;
DROP TABLE IF EXISTS user_favorites;
DROP TABLE IF EXISTS song_plays;
DROP TABLE IF EXISTS songs;

-- +goose StatementEnd 