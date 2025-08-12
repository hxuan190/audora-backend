// pkg/queue/celery.go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"music-app-backend/pkg/redis"
	"time"

	"github.com/google/uuid"
)

const (
	CeleryDefaultQueue = "audio_processing" // Celery queue name
	CeleryResultPrefix = "celery-task-meta-" // Celery result key prefix
)

type CeleryClient struct {
	redisClient *redis.Client
	broker      string // Redis broker URL
}

// CeleryTask represents a Celery task in the format Celery expects
type CeleryTask struct {
	ID        string                 `json:"id"`
	Task      string                 `json:"task"`      // Task name (e.g., "audio_processor.process_track")
	Args      []interface{}          `json:"args"`      // Positional arguments
	Kwargs    map[string]interface{} `json:"kwargs"`    // Keyword arguments
	Retries   int                    `json:"retries"`
	ETA       *time.Time             `json:"eta,omitempty"`       // Estimated time of arrival
	Expires   *time.Time             `json:"expires,omitempty"`   // Task expiration
	UTCOffset int                    `json:"utc,omitempty"`       // UTC offset
	Callbacks []string               `json:"callbacks,omitempty"` // Success callbacks
	Errbacks  []string               `json:"errbacks,omitempty"`  // Error callbacks
	TimeLimit []int                  `json:"timelimit,omitempty"` // [soft, hard] time limits
	Taskset   string                 `json:"taskset,omitempty"`   // Task group ID
}

// AudioProcessingTask represents the specific task data for audio processing
type AudioProcessingTask struct {
	SongID           uint64                    `json:"song_id"`
	ArtistID         uint64                    `json:"artist_id"`
	SourceBucket     string                    `json:"source_bucket"`
	SourceObjectPath string                    `json:"source_object_path"`
	DestBucket       string                    `json:"dest_bucket"`
	ProcessingConfig AudioProcessingConfig     `json:"processing_config"`
	Metadata         AudioProcessingMetadata   `json:"metadata"`
	CallbackURL      string                    `json:"callback_url,omitempty"` // Optional webhook callback
}

type AudioProcessingConfig struct {
	TargetLUFS           float64  `json:"target_lufs"`            // -14.0 default
	GenerateFormats      []string `json:"generate_formats"`       // ["mp3_320", "flac_cd"]
	QualityEnhancement   bool     `json:"quality_enhancement"`    // Apply noise reduction, etc.
	PreserveDynamicRange bool     `json:"preserve_dynamic_range"` // Avoid over-compression
	ProcessingIntensity  string   `json:"processing_intensity"`   // "conservative", "standard", "aggressive"
	ValidateOnly         bool     `json:"validate_only"`          // Only validate, don't process
}

type AudioProcessingMetadata struct {
	OriginalFilename string            `json:"original_filename"`
	FileSize         int64             `json:"file_size"`
	ContentType      string            `json:"content_type"`
	UploadSessionID  string            `json:"upload_session_id"`
	Title            string            `json:"title"`
	GenreID          *uint64           `json:"genre_id,omitempty"`
	MoodID           *uint64           `json:"mood_id,omitempty"`
	Description      string            `json:"description,omitempty"`
	AdditionalData   map[string]string `json:"additional_data,omitempty"`
}

// CeleryTaskResult represents the result from a Celery task
type CeleryTaskResult struct {
	TaskID    string      `json:"task_id"`
	Status    string      `json:"status"`    // "PENDING", "STARTED", "SUCCESS", "FAILURE", "RETRY", "REVOKED"
	Result    interface{} `json:"result"`    // Task result data
	Traceback string      `json:"traceback"` // Error traceback if failed
	Children  []string    `json:"children"`  // Child task IDs
	DateDone  *time.Time  `json:"date_done"` // Completion timestamp
}

// AudioProcessingResult represents the result from audio processing
type AudioProcessingResult struct {
	SongID           uint64                 `json:"song_id"`
	Success          bool                   `json:"success"`
	ProcessedFormats []ProcessedAudioFormat `json:"processed_formats"`
	AudioAnalysis    AudioAnalysis          `json:"audio_analysis"`
	QualityScore     float64                `json:"quality_score"`
	ProcessingTime   float64                `json:"processing_time_seconds"`
	Warnings         []string               `json:"warnings,omitempty"`
	Error            string                 `json:"error,omitempty"`
	MasteredForAudora bool                  `json:"mastered_for_audora"`
}

type ProcessedAudioFormat struct {
	Format       string  `json:"format"`        // "mp3_320", "flac_cd", "flac_hires"
	ObjectPath   string  `json:"object_path"`   // Path in processed-tracks bucket
	FileSize     int64   `json:"file_size"`     // Size in bytes
	Bitrate      int     `json:"bitrate"`       // Actual bitrate
	SampleRate   int     `json:"sample_rate"`   // Sample rate in Hz
	BitDepth     int     `json:"bit_depth"`     // Bit depth (for lossless)
	Duration     float64 `json:"duration"`      // Duration in seconds
	QualityScore float64 `json:"quality_score"` // Format-specific quality score
}

type AudioAnalysis struct {
	OriginalFormat    string  `json:"original_format"`
	OriginalBitrate   int     `json:"original_bitrate"`
	OriginalSampleRate int    `json:"original_sample_rate"`
	OriginalBitDepth  int     `json:"original_bit_depth"`
	Duration          float64 `json:"duration"`
	OriginalLUFS      float64 `json:"original_lufs"`
	ProcessedLUFS     float64 `json:"processed_lufs"`
	DynamicRange      float64 `json:"dynamic_range"`      // DR measurement
	PeakLevel         float64 `json:"peak_level"`         // Peak dBFS
	TruePeak          float64 `json:"true_peak"`          // True peak dBTP
	SpectralCentroid  float64 `json:"spectral_centroid"`  // Frequency analysis
	THDPlusN          float64 `json:"thd_plus_n"`         // Total harmonic distortion + noise
	StereoWidth       float64 `json:"stereo_width"`       // Stereo field width
	HasClipping       bool    `json:"has_clipping"`       // Digital clipping detected
	HasArtifacts      bool    `json:"has_artifacts"`      // Processing artifacts detected
	QualityGrade      string  `json:"quality_grade"`      // "studio", "mastered", "good", "needs_improvement"
}

func NewCeleryClient(redisClient *redis.Client) *CeleryClient {
	return &CeleryClient{
		redisClient: redisClient,
		broker:      fmt.Sprintf("redis://%s", redisClient), // Will be properly formatted
	}
}

// SubmitAudioProcessingTask submits an audio processing task to Celery
func (c *CeleryClient) SubmitAudioProcessingTask(ctx context.Context, taskData *AudioProcessingTask) (string, error) {
	taskID := uuid.New().String()

	// Create Celery task in the format Celery expects
	celeryTask := &CeleryTask{
		ID:   taskID,
		Task: "audio_processor.process_track", // Python task name
		Args: []interface{}{},                  // Use kwargs instead of args for clarity
		Kwargs: map[string]interface{}{
			"song_id":             taskData.SongID,
			"artist_id":           taskData.ArtistID,
			"source_bucket":       taskData.SourceBucket,
			"source_object_path":  taskData.SourceObjectPath,
			"dest_bucket":         taskData.DestBucket,
			"processing_config":   taskData.ProcessingConfig,
			"metadata":            taskData.Metadata,
			"callback_url":        taskData.CallbackURL,
		},
		Retries: 3, // Allow up to 3 retries
	}

	// Serialize task to JSON
	taskJSON, err := json.Marshal(celeryTask)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Celery task: %w", err)
	}

	// Create Celery message format
	celeryMessage := map[string]interface{}{
		"body":         string(taskJSON),
		"content-type": "application/json",
		"content-encoding": "utf-8",
		"headers": map[string]interface{}{
			"id":           taskID,
			"task":         celeryTask.Task,
			"lang":         "py",
			"root_id":      taskID,
			"parent_id":    nil,
			"group":        nil,
			"meth":         "py",
			"shadow":       nil,
			"eta":          nil,
			"expires":      nil,
			"retries":      0,
			"timelimit":    []interface{}{nil, nil},
			"argsrepr":     "[]",
			"kwargsrepr":   fmt.Sprintf("%v", celeryTask.Kwargs),
			"origin":       "audora-api",
		},
		"properties": map[string]interface{}{
			"correlation_id": taskID,
			"reply_to":       uuid.New().String(),
			"delivery_mode":  2,
			"delivery_info": map[string]interface{}{
				"exchange":    "",
				"routing_key": CeleryDefaultQueue,
			},
		},
	}

	// Serialize the complete message
	messageJSON, err := json.Marshal(celeryMessage)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Celery message: %w", err)
	}

	// Push to Redis list (Celery queue)
	err = c.redisClient.LPush(ctx, CeleryDefaultQueue, string(messageJSON))
	if err != nil {
		return "", fmt.Errorf("failed to enqueue Celery task: %w", err)
	}

	return taskID, nil
}

// GetTaskResult retrieves the result of a Celery task
func (c *CeleryClient) GetTaskResult(ctx context.Context, taskID string) (*CeleryTaskResult, error) {
	resultKey := CeleryResultPrefix + taskID
	
	resultJSON, err := c.redisClient.Get(ctx, resultKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			// Task result not yet available
			return &CeleryTaskResult{
				TaskID: taskID,
				Status: "PENDING",
			}, nil
		}
		return nil, fmt.Errorf("failed to get task result: %w", err)
	}

	var result CeleryTaskResult
	err = json.Unmarshal([]byte(resultJSON), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task result: %w", err)
	}

	return &result, nil
}

// GetAudioProcessingResult retrieves and parses audio processing result
func (c *CeleryClient) GetAudioProcessingResult(ctx context.Context, taskID string) (*AudioProcessingResult, error) {
	celeryResult, err := c.GetTaskResult(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if celeryResult.Status != "SUCCESS" {
		return &AudioProcessingResult{
			Success: false,
			Error:   fmt.Sprintf("Task status: %s", celeryResult.Status),
		}, nil
	}

	// Parse the result data as AudioProcessingResult
	resultBytes, err := json.Marshal(celeryResult.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result data: %w", err)
	}

	var audioResult AudioProcessingResult
	err = json.Unmarshal(resultBytes, &audioResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal audio processing result: %w", err)
	}

	return &audioResult, nil
}

// RevokeTask cancels a pending or running Celery task
func (c *CeleryClient) RevokeTask(ctx context.Context, taskID string, terminate bool) error {
	revokeData := map[string]interface{}{
		"method":    "revoke",
		"arguments": map[string]interface{}{
			"task_id":   taskID,
			"terminate": terminate,
			"signal":    "SIGTERM",
		},
	}

	revokeJSON, err := json.Marshal(revokeData)
	if err != nil {
		return fmt.Errorf("failed to marshal revoke data: %w", err)
	}

	// Send revoke command to Celery management queue
	return c.redisClient.LPush(ctx, "celeryev.management", string(revokeJSON))
}

// GetQueueStats returns statistics about Celery queues
func (c *CeleryClient) GetQueueStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Get main queue length
	queueLen, err := c.redisClient.LLen(ctx, CeleryDefaultQueue)
	if err != nil {
		return nil, err
	}
	stats["audio_processing_queue"] = queueLen

	// Get active tasks count (approximate)
	// This would need to be implemented based on your Celery monitoring setup
	// stats["active_tasks"] = activeCount

	return stats, nil
}

// WaitForResult waits for a task to complete with timeout
func (c *CeleryClient) WaitForResult(ctx context.Context, taskID string, timeout time.Duration) (*CeleryTaskResult, error) {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		result, err := c.GetTaskResult(ctx, taskID)
		if err != nil {
			return nil, err
		}

		if result.Status != "PENDING" && result.Status != "STARTED" {
			return result, nil
		}

		// Wait a bit before checking again
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("task did not complete within timeout: %v", timeout)
}

// CreateProcessingTaskForSong creates a processing task with sensible defaults
func (c *CeleryClient) CreateProcessingTaskForSong(
	songID, artistID uint64,
	sourceObjectPath string,
	metadata AudioProcessingMetadata,
	callbackURL string,
) *AudioProcessingTask {
	return &AudioProcessingTask{
		SongID:           songID,
		ArtistID:         artistID,
		SourceBucket:     "audora-tracks",
		SourceObjectPath: sourceObjectPath,
		DestBucket:       "processed-tracks",
		CallbackURL:      callbackURL,
		ProcessingConfig: AudioProcessingConfig{
			TargetLUFS:           -14.0,
			GenerateFormats:      []string{"mp3_320", "flac_cd"},
			QualityEnhancement:   true,
			PreserveDynamicRange: true,
			ProcessingIntensity:  "standard",
			ValidateOnly:         false,
		},
		Metadata: metadata,
	}
}