-- Music App MVP Database Schema
-- Designed for PostgreSQL with Kratos authentication integration

-- Users table (linked to Kratos identities)
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    kratos_identity_id UUID NOT NULL UNIQUE, -- Links to Kratos identity
    email VARCHAR(255) NOT NULL UNIQUE,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('artist', 'listener', 'admin')),
    display_name VARCHAR(100),
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Artists table (extends user for artist-specific data)
CREATE TABLE artists (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL, -- No FK reference
    artist_name VARCHAR(150) NOT NULL,
    bio TEXT,
    profile_image_url TEXT,
    banner_image_url TEXT,
    website_url TEXT,
    spotify_url TEXT,
    instagram_url TEXT,
    twitter_url TEXT,
    youtube_url TEXT,
    is_verified BOOLEAN DEFAULT false,
    verification_requested_at TIMESTAMP WITH TIME ZONE,
    total_plays BIGINT DEFAULT 0,
    total_earnings DECIMAL(10,2) DEFAULT 0.00,
    follower_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id)
);

-- Content tiers enum
CREATE TYPE content_tier AS ENUM ('public_discovery', 'fan_exclusives', 'collaboration_hub', 'personal_archive');

-- Genres and moods lookup tables
CREATE TABLE genres (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE moods (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    color_hex VARCHAR(7), -- For UI theming
    is_active BOOLEAN DEFAULT true
);

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
    played_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_song_plays_song_id (song_id),
    INDEX idx_song_plays_user_id (user_id),
    INDEX idx_song_plays_played_at (played_at)
);

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
    is_active BOOLEAN DEFAULT true,
    
    INDEX idx_listening_sessions_artist_id (artist_id),
    INDEX idx_listening_sessions_active (is_active, last_heartbeat)
);

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

-- Analytics aggregation tables for performance
CREATE TABLE daily_artist_stats (
    id UUID PRIMARY KEY NOT NULL,
    artist_id UUID NOT NULL, -- No FK reference
    date DATE NOT NULL,
    total_plays INTEGER DEFAULT 0,
    unique_listeners INTEGER DEFAULT 0,
    total_duration_played INTEGER DEFAULT 0,
    new_followers INTEGER DEFAULT 0,
    tips_received_cents INTEGER DEFAULT 0,
    tip_count INTEGER DEFAULT 0,
    
    UNIQUE(artist_id, date)
);

CREATE TABLE daily_song_stats (
    id UUID PRIMARY KEY NOT NULL,
    song_id UUID NOT NULL, -- No FK reference
    date DATE NOT NULL,
    play_count INTEGER DEFAULT 0,
    unique_listeners INTEGER DEFAULT 0,
    completion_rate DECIMAL(5,2) DEFAULT 0.00,
    avg_duration_played INTEGER DEFAULT 0,
    skip_rate DECIMAL(5,2) DEFAULT 0.00,
    
    UNIQUE(song_id, date)
);

-- Indexes for performance
CREATE INDEX idx_users_kratos_identity ON users(kratos_identity_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_type ON users(user_type);

CREATE INDEX idx_songs_artist_id ON songs(artist_id);
CREATE INDEX idx_songs_genre_id ON songs(genre_id);
CREATE INDEX idx_songs_mood_id ON songs(mood_id);
CREATE INDEX idx_songs_tier ON songs(tier);
CREATE INDEX idx_songs_active ON songs(is_active);
CREATE INDEX idx_songs_created_at ON songs(created_at);

CREATE INDEX idx_tips_artist_id ON tips(to_artist_id);
CREATE INDEX idx_tips_user_id ON tips(from_user_id);
CREATE INDEX idx_tips_status ON tips(status);
CREATE INDEX idx_tips_created_at ON tips(created_at);

-- Triggers for updating aggregate counts
CREATE OR REPLACE FUNCTION update_artist_totals()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Update play count and earnings
        UPDATE artists 
        SET total_plays = total_plays + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = (SELECT artist_id FROM songs WHERE id = NEW.song_id);
        
        -- Update song play count
        UPDATE songs 
        SET play_count = play_count + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.song_id;
        
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_play_counts
    AFTER INSERT ON song_plays
    FOR EACH ROW
    EXECUTE FUNCTION update_artist_totals();

-- Function to update follower counts
CREATE OR REPLACE FUNCTION update_follower_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE artists 
        SET follower_count = follower_count + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.artist_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE artists 
        SET follower_count = follower_count - 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.artist_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_follower_counts
    AFTER INSERT OR DELETE ON artist_followers
    FOR EACH ROW
    EXECUTE FUNCTION update_follower_count();

-- Function to update tip totals
CREATE OR REPLACE FUNCTION update_tip_totals()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status != 'completed') THEN
        -- Update artist earnings
        UPDATE artists 
        SET total_earnings = total_earnings + (NEW.artist_payout_cents / 100.0),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.to_artist_id;
        
        -- Update song tip stats if specific song
        IF NEW.song_id IS NOT NULL THEN
            UPDATE songs 
            SET tip_count = tip_count + 1,
                total_tips = total_tips + (NEW.artist_payout_cents / 100.0),
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.song_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tip_totals
    AFTER UPDATE ON tips
    FOR EACH ROW
    EXECUTE FUNCTION update_tip_totals();

-- Insert initial data
INSERT INTO genres (name, description) VALUES
('Electronic', 'Electronic music and EDM'),
('Hip Hop', 'Hip hop and rap music'),
('Rock', 'Rock and alternative music'),
('Pop', 'Popular music'),
('Jazz', 'Jazz and fusion'),
('Classical', 'Classical and orchestral'),
('R&B', 'R&B and soul music'),
('Country', 'Country and folk'),
('Indie', 'Independent and alternative'),
('World', 'World and ethnic music');

INSERT INTO moods (name, description, color_hex) VALUES
('Focus', 'Music for concentration and productivity', '#4A90E2'),
('Workout', 'High energy music for exercise', '#E74C3C'),
('Chill', 'Relaxed and mellow vibes', '#2ECC71'),
('Morning', 'Uplifting music to start the day', '#F39C12'),
('Evening', 'Calm music for winding down', '#8E44AD'),
('Party', 'Upbeat music for celebrations', '#E91E63'),
('Study', 'Ambient music for learning', '#607D8B'),
('Sleep', 'Peaceful music for rest', '#34495E'),
('Drive', 'Music for road trips', '#FF5722'),
('Romance', 'Music for romantic moments', '#E91E63');

-- Create initial mood-based playlists
INSERT INTO playlists (name, description, playlist_type, mood_id) VALUES
('Focus Flow', 'Instrumental tracks perfect for deep work and concentration', 'curated', 1),
('Workout Warriors', 'High-energy beats to power your workout', 'curated', 2),
('Chill Vibes', 'Laid-back tracks for relaxation', 'curated', 3),
('Morning Motivation', 'Uplifting songs to start your day right', 'curated', 4),
('Evening Wind Down', 'Mellow tunes for a peaceful evening', 'curated', 5);