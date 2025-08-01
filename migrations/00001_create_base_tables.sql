-- +goose Up
-- +goose StatementBegin

-- Users table (linked to Kratos identities)
CREATE TABLE users (
    id BIGINT PRIMARY KEY NOT NULL,
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
    id BIGINT PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL, -- No FK reference
    artist_name VARCHAR(150) NOT NULL,
    bio TEXT DEFAULT NULL,
    profile_image_url TEXT DEFAULT NULL,
    banner_image_url TEXT DEFAULT NULL,
    website_url TEXT DEFAULT NULL,
    spotify_url TEXT DEFAULT NULL,
    instagram_url TEXT DEFAULT NULL,
    twitter_url TEXT DEFAULT NULL,
    youtube_url TEXT DEFAULT NULL,
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
    id BIGINT PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE moods (
    id BIGINT PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    color_hex VARCHAR(7), -- For UI theming
    is_active BOOLEAN DEFAULT true
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS moods;
DROP TABLE IF EXISTS genres;
DROP TYPE IF EXISTS content_tier;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd 