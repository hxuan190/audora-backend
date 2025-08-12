package model

import (
	"music-app-backend/pkg/model"
)

// ProcessedAudioFormat represents a processed version of a song in a specific format
type ProcessedAudioFormat struct {
	model.BaseModel
	SongID           uint64  `json:"song_id" gorm:"not null;index"`
	Format           string  `json:"format" gorm:"not null;size:50;index"`   // mp3_320, flac_cd, flac_hires
	ObjectPath       string  `json:"object_path" gorm:"not null"`            // Path in processed-tracks bucket
	FileSize         int64   `json:"file_size" gorm:"not null"`              // Size in bytes
	Bitrate          *int    `json:"bitrate"`                                // Actual bitrate
	SampleRate       *int    `json:"sample_rate"`                            // Sample rate in Hz
	BitDepth         *int    `json:"bit_depth"`                              // Bit depth (for lossless)
	Duration         float64 `json:"duration" gorm:"type:decimal(8,3)"`      // Duration in seconds
	QualityScore     float64 `json:"quality_score" gorm:"type:decimal(5,3)"` // Format-specific quality score (0-1)
	ProcessingTaskID string  `json:"processing_task_id" gorm:"size:100"`     // Reference to Celery task
}

// TableName returns the table name for ProcessedAudioFormat
func (ProcessedAudioFormat) TableName() string {
	return "processed_audio_formats"
}

// AudioAnalysis represents detailed audio analysis results for a song
type AudioAnalysis struct {
	model.BaseModel
	SongID             uint64   `json:"song_id" gorm:"not null;uniqueIndex"` // One analysis per song
	OriginalFormat     string   `json:"original_format" gorm:"size:50"`
	OriginalBitrate    *int     `json:"original_bitrate"`
	OriginalSampleRate *int     `json:"original_sample_rate"`
	OriginalBitDepth   *int     `json:"original_bit_depth"`
	Duration           float64  `json:"duration" gorm:"type:decimal(8,3)"`          // Duration in seconds
	OriginalLUFS       *float64 `json:"original_lufs" gorm:"type:decimal(6,2)"`     // Original LUFS measurement
	ProcessedLUFS      *float64 `json:"processed_lufs" gorm:"type:decimal(6,2)"`    // Processed LUFS measurement
	DynamicRange       *float64 `json:"dynamic_range" gorm:"type:decimal(6,2)"`     // DR measurement
	PeakLevel          *float64 `json:"peak_level" gorm:"type:decimal(6,2)"`        // Peak dBFS
	TruePeak           *float64 `json:"true_peak" gorm:"type:decimal(6,2)"`         // True peak dBTP
	SpectralCentroid   *float64 `json:"spectral_centroid" gorm:"type:decimal(8,2)"` // Frequency analysis
	THDPlusN           *float64 `json:"thd_plus_n" gorm:"type:decimal(8,4)"`        // Total harmonic distortion + noise
	StereoWidth        *float64 `json:"stereo_width" gorm:"type:decimal(5,3)"`      // Stereo field width
	HasClipping        bool     `json:"has_clipping" gorm:"default:false"`          // Digital clipping detected
	HasArtifacts       bool     `json:"has_artifacts" gorm:"default:false"`         // Processing artifacts detected
	QualityGrade       string   `json:"quality_grade" gorm:"size:50;index"`         // studio, mastered, good, needs_improvement
	QualityScore       float64  `json:"quality_score" gorm:"type:decimal(5,3)"`     // Overall quality score (0-1)
	ProcessingTime     float64  `json:"processing_time" gorm:"type:decimal(8,3)"`   // Processing time in seconds
	Warnings           []string `json:"warnings" gorm:"type:text[]"`                // Array of warning messages
	ProcessingTaskID   string   `json:"processing_task_id" gorm:"size:100"`         // Reference to Celery task
}

// TableName returns the table name for AudioAnalysis
func (AudioAnalysis) TableName() string {
	return "audio_analysis"
}
