-- +goose Up
-- +goose StatementBegin

-- Processed audio formats table to store different quality versions of songs
CREATE TABLE processed_audio_formats (
    id BIGINT PRIMARY KEY NOT NULL,
    song_id BIGINT NOT NULL, -- No FK reference
    format VARCHAR(50) NOT NULL, -- mp3_320, flac_cd, flac_hires, etc.
    object_path TEXT NOT NULL, -- Path in processed-tracks bucket
    file_size BIGINT NOT NULL, -- Size in bytes
    bitrate INTEGER, -- Actual bitrate
    sample_rate INTEGER, -- Sample rate in Hz
    bit_depth INTEGER, -- Bit depth (for lossless)
    duration DECIMAL(8,3), -- Duration in seconds with precision
    quality_score DECIMAL(5,3), -- Format-specific quality score (0-1)
    processing_task_id VARCHAR(100), -- Reference to Celery task
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(song_id, format)
);

-- Audio analysis results table to store detailed analysis
CREATE TABLE audio_analysis (
    id BIGINT PRIMARY KEY NOT NULL,
    song_id BIGINT NOT NULL UNIQUE, -- No FK reference, one analysis per song
    original_format VARCHAR(50),
    original_bitrate INTEGER,
    original_sample_rate INTEGER,
    original_bit_depth INTEGER,
    duration DECIMAL(8,3), -- Duration in seconds
    original_lufs DECIMAL(6,2), -- Original LUFS measurement
    processed_lufs DECIMAL(6,2), -- Processed LUFS measurement
    dynamic_range DECIMAL(6,2), -- DR measurement
    peak_level DECIMAL(6,2), -- Peak dBFS
    true_peak DECIMAL(6,2), -- True peak dBTP
    spectral_centroid DECIMAL(8,2), -- Frequency analysis
    thd_plus_n DECIMAL(8,4), -- Total harmonic distortion + noise
    stereo_width DECIMAL(5,3), -- Stereo field width
    has_clipping BOOLEAN DEFAULT false, -- Digital clipping detected
    has_artifacts BOOLEAN DEFAULT false, -- Processing artifacts detected
    quality_grade VARCHAR(50), -- studio, mastered, good, needs_improvement
    quality_score DECIMAL(5,3), -- Overall quality score (0-1)
    processing_time DECIMAL(8,3), -- Processing time in seconds
    warnings TEXT[], -- Array of warning messages
    processing_task_id VARCHAR(100), -- Reference to Celery task
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_processed_audio_formats_song_id ON processed_audio_formats(song_id);
CREATE INDEX idx_processed_audio_formats_format ON processed_audio_formats(format);
CREATE INDEX idx_audio_analysis_song_id ON audio_analysis(song_id);
CREATE INDEX idx_audio_analysis_quality_grade ON audio_analysis(quality_grade);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_audio_analysis_quality_grade;
DROP INDEX IF EXISTS idx_audio_analysis_song_id;
DROP INDEX IF EXISTS idx_processed_audio_formats_format;
DROP INDEX IF EXISTS idx_processed_audio_formats_song_id;
DROP TABLE IF EXISTS audio_analysis;
DROP TABLE IF EXISTS processed_audio_formats;

-- +goose StatementEnd
