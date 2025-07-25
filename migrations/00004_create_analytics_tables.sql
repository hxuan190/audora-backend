-- +goose Up
-- +goose StatementBegin

-- Analytics aggregation tables for performance
CREATE TABLE daily_artist_stats (
    id BIGINT PRIMARY KEY NOT NULL,
    artist_id BIGINT NOT NULL, -- No FK reference
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
    id BIGINT PRIMARY KEY NOT NULL,
    song_id BIGINT NOT NULL, -- No FK reference
    date DATE NOT NULL,
    play_count INTEGER DEFAULT 0,
    unique_listeners INTEGER DEFAULT 0,
    completion_rate DECIMAL(5,2) DEFAULT 0.00,
    avg_duration_played INTEGER DEFAULT 0,
    skip_rate DECIMAL(5,2) DEFAULT 0.00,
    
    UNIQUE(song_id, date)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS daily_song_stats;
DROP TABLE IF EXISTS daily_artist_stats;

-- +goose StatementEnd 