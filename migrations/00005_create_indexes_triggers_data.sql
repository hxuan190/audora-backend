-- +goose Up
-- +goose StatementBegin

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_kratos_identity ON users(kratos_identity_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_type ON users(user_type);

CREATE INDEX IF NOT EXISTS idx_songs_artist_id ON songs(artist_id);
CREATE INDEX IF NOT EXISTS idx_songs_genre_id ON songs(genre_id);
CREATE INDEX IF NOT EXISTS idx_songs_mood_id ON songs(mood_id);
CREATE INDEX IF NOT EXISTS idx_songs_tier ON songs(tier);
CREATE INDEX IF NOT EXISTS idx_songs_active ON songs(is_active);
CREATE INDEX IF NOT EXISTS idx_songs_created_at ON songs(created_at);

CREATE INDEX IF NOT EXISTS idx_tips_artist_id ON tips(to_artist_id);
CREATE INDEX IF NOT EXISTS idx_tips_user_id ON tips(from_user_id);
CREATE INDEX IF NOT EXISTS idx_tips_status ON tips(status);
CREATE INDEX IF NOT EXISTS idx_tips_created_at ON tips(created_at);

-- Triggers for updating aggregate counts
CREATE OR REPLACE FUNCTION update_artist_totals()
RETURNS TRIGGER AS $$
DECLARE
    artist_id BIGINT;
BEGIN
    SELECT artist_id INTO artist_id FROM songs WHERE id = NEW.song_id;

    -- Update song
    UPDATE songs
    SET play_count = play_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.song_id;

    -- Update artist
    UPDATE artists
    SET total_plays = total_plays + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = artist_id;

    RETURN NEW;
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
        SET follower_count = GREATEST(follower_count - 1, 0),
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
    IF NEW.status = 'completed' AND (OLD.status IS DISTINCT FROM 'completed') THEN        -- Update artist earnings
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
INSERT INTO genres (id, name, description) VALUES
(1000000000000000001, 'Electronic', 'Electronic music and EDM'),
(1000000000000000002, 'Hip Hop', 'Hip hop and rap music'),
(1000000000000000003, 'Rock', 'Rock and alternative music'),
(1000000000000000004, 'Pop', 'Popular music'),
(1000000000000000005, 'Jazz', 'Jazz and fusion'),
(1000000000000000006, 'Classical', 'Classical and orchestral'),
(1000000000000000007, 'R&B', 'R&B and soul music'),
(1000000000000000008, 'Country', 'Country and folk'),
(1000000000000000009, 'Indie', 'Independent and alternative'),
(1000000000000000010, 'World', 'World and ethnic music');

INSERT INTO moods (id, name, description, color_hex) VALUES
(1000000000000000011, 'Focus', 'Music for concentration and productivity', '#4A90E2'),
(1000000000000000012, 'Workout', 'High energy music for exercise', '#E74C3C'),
(1000000000000000013, 'Chill', 'Relaxed and mellow vibes', '#2ECC71'),
(1000000000000000014, 'Morning', 'Uplifting music to start the day', '#F39C12'),
(1000000000000000015, 'Evening', 'Calm music for winding down', '#8E44AD'),
(1000000000000000016, 'Party', 'Upbeat music for celebrations', '#E91E63'),
(1000000000000000017, 'Study', 'Ambient music for learning', '#607D8B'),
(1000000000000000018, 'Sleep', 'Peaceful music for rest', '#34495E'),
(1000000000000000019, 'Drive', 'Music for road trips', '#FF5722'),
(1000000000000000020, 'Romance', 'Music for romantic moments', '#E91E63');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove initial data
DELETE FROM moods;
DELETE FROM genres;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_tip_totals ON tips;
DROP TRIGGER IF EXISTS trigger_update_follower_counts ON artist_followers;
DROP TRIGGER IF EXISTS trigger_update_play_counts ON song_plays;

-- Drop functions
DROP FUNCTION IF EXISTS update_tip_totals();
DROP FUNCTION IF EXISTS update_follower_count();
DROP FUNCTION IF EXISTS update_artist_totals();

-- Drop indexes
DROP INDEX IF EXISTS idx_tips_created_at;
DROP INDEX IF EXISTS idx_tips_status;
DROP INDEX IF EXISTS idx_tips_user_id;
DROP INDEX IF EXISTS idx_tips_artist_id;
DROP INDEX IF EXISTS idx_songs_created_at;
DROP INDEX IF EXISTS idx_songs_active;
DROP INDEX IF EXISTS idx_songs_tier;
DROP INDEX IF EXISTS idx_songs_mood_id;
DROP INDEX IF EXISTS idx_songs_genre_id;
DROP INDEX IF EXISTS idx_songs_artist_id;
DROP INDEX IF EXISTS idx_users_type;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_kratos_identity;

-- +goose StatementEnd 