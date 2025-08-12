// internal/music/adapters/http/upload_handler.go - Updated with Celery integration
package http

import (
	"fmt"
	"music-app-backend/internal/music/application"
	model "music-app-backend/internal/music/domain"
	appError "music-app-backend/pkg/error"
	jsonResponse "music-app-backend/pkg/json"
	"music-app-backend/pkg/queue"
	"music-app-backend/pkg/redis"
	"music-app-backend/pkg/storage"
	"path/filepath"
	"strings"
	"time"

	baseModel "music-app-backend/pkg/model"

	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	musicService   *application.MusicService
	storageService *storage.MinIOService
	celeryClient   *queue.CeleryClient
	generator      *goflakeid.Generator
}

type InitiateUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	ArtistID    uint64 `json:"artist_id" binding:"required"`
}

type InitiateUploadResponse struct {
	UploadID     string              `json:"upload_id"`
	UploadURL    string              `json:"upload_url"`
	ExpiresAt    time.Time           `json:"expires_at"`
	MaxFileSize  int64               `json:"max_file_size"`
	Instructions *UploadInstructions `json:"instructions"`
}

type UploadInstructions struct {
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	CallbackURL string            `json:"callback_url"`
}

type CompleteUploadRequest struct {
	// Upload confirmation data
	UploadID   string `json:"upload_id" binding:"required"`
	FileURL    string `json:"file_url" binding:"required"`
	ActualSize int64  `json:"actual_size" binding:"required"`
	ETag       string `json:"etag"`

	// Required song metadata
	Title   string  `json:"title" binding:"required,min=1,max=200"`
	GenreID *uint64 `json:"genre_id" binding:"required"`
	MoodID  *uint64 `json:"mood_id" binding:"required"`

	// Optional song metadata
	Description      string                       `json:"description"`
	ProcessingConfig *AudioProcessingConfigRequest `json:"processing_config,omitempty"`
}

type AudioProcessingConfigRequest struct {
	TargetLUFS           *float64 `json:"target_lufs,omitempty"`            // -14.0 default
	GenerateFormats      []string `json:"generate_formats,omitempty"`       // ["mp3_320", "flac_cd", "flac_hires"]
	QualityEnhancement   *bool    `json:"quality_enhancement,omitempty"`    // Apply additional processing
	PreserveDynamicRange *bool    `json:"preserve_dynamic_range,omitempty"` // Avoid over-compression
	ProcessingIntensity  string   `json:"processing_intensity,omitempty"`   // "conservative", "standard", "aggressive"
}

type CompleteUploadResponse struct {
	SongID              uint64         `json:"song_id"`
	Status              string         `json:"status"`
	Message             string         `json:"message"`
	ProcessingTaskID    string         `json:"processing_task_id"`
	EstimatedCompletion time.Time      `json:"estimated_completion"`
	UploadSummary       *UploadSummary `json:"upload_summary"`
	NextSteps           []string       `json:"next_steps"`
	TrackingURL         string         `json:"tracking_url"`
}

type UploadSummary struct {
	FileSize   int64     `json:"file_size"`
	Format     string    `json:"format"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type ProcessingStatusResponse struct {
	SongID      uint64                       `json:"song_id"`
	TaskID      string                       `json:"task_id"`
	Status      string                       `json:"status"` // "PENDING", "STARTED", "SUCCESS", "FAILURE"
	Progress    *ProcessingProgress          `json:"progress,omitempty"`
	Result      *queue.AudioProcessingResult `json:"result,omitempty"`
	Error       string                       `json:"error,omitempty"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
}

type ProcessingProgress struct {
	Stage       string  `json:"stage"`       // "downloading", "validating", "normalizing", "transcoding", "uploading"
	Percentage  float64 `json:"percentage"`  // 0.0 to 100.0
	CurrentStep string  `json:"current_step"`
	ETA         *time.Time `json:"eta,omitempty"`
}

func NewMusicHandler(
	musicService *application.MusicService, 
	storageService *storage.MinIOService, 
	redisClient *redis.Client,
	generator *goflakeid.Generator,
) *MusicHandler {
	celeryClient := queue.NewCeleryClient(redisClient)

	return &MusicHandler{
		musicService:   musicService,
		storageService: storageService,
		celeryClient:   celeryClient,
		generator:      generator,
	}
}

func (h *MusicHandler) HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if appErr, ok := appError.GetAppError(err); ok {
		jsonResponse.ResponseJSON(c, appErr.StatusCode, appErr.Message, appErr.Data)
		return true
	}

	jsonResponse.ResponseInternalError(c, err)
	return true
}

func (h *MusicHandler) InitiateUpload(c *gin.Context) {
	request := &InitiateUploadRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonResponse.ResponseBadRequest(c, "Invalid request: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	// Validate audio format
	if !h.isValidAudioFormat(request.Filename) {
		jsonResponse.ResponseBadRequest(c, "Unsupported file format. Supported: FLAC, WAV, AIFF, MP3")
		return
	}

	// Validate file size
	maxSize := int64(600 * 1024 * 1024) // 600MB
	if request.FileSize <= 0 || request.FileSize > maxSize {
		jsonResponse.ResponseBadRequest(c, "File size must be between 1 byte and 600MB")
		return
	}

	// Generate upload session
	uploadID := h.generateUploadID(request.ArtistID, request.Filename)
	objectPath := h.storageService.GenerateUploadPath(request.ArtistID, request.Filename)

	// Get presigned upload URL
	presignedURL, err := h.storageService.GetPresignedUploadURL(
		c.Request.Context(), 
		objectPath, 
		storage.BucketTypeTracks, 
		15*time.Minute,
	)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to generate presigned URL: %v", err))
		return
	}

	// Create upload session
	uploadSession := &model.UploadSession{
		ID:         uploadID,
		ArtistID:   request.ArtistID,
		UserID:     userID.(uint64),
		Filename:   request.Filename,
		FileSize:   request.FileSize,
		ObjectPath: objectPath,
		Status:     "initiated",
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}

	_, err = h.musicService.CreateUploadSession(c.Request.Context(), uploadSession)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to create upload session: %v", err))
		return
	}

	response := &InitiateUploadResponse{
		UploadID:    uploadID,
		UploadURL:   presignedURL.URL,
		ExpiresAt:   presignedURL.ExpiresAt,
		MaxFileSize: maxSize,
		Instructions: &UploadInstructions{
			Method:      presignedURL.Method,
			Headers:     presignedURL.Headers,
			CallbackURL: "/api/v1/upload/complete",
		},
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *MusicHandler) CompleteUpload(c *gin.Context) {
	request := &CompleteUploadRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonResponse.ResponseBadRequest(c, "Invalid request: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	// Verify upload session
	uploadSession, err := h.musicService.GetUploadSession(c.Request.Context(), request.UploadID)
	if h.HandleError(c, err) {
		return
	}

	if uploadSession.UserID != userID.(uint64) {
		jsonResponse.ResponseForbidden(c)
		return
	}

	// Verify file was uploaded correctly
	objectPath := h.extractObjectPathFromURL(request.FileURL)
	fileInfo, err := h.storageService.GetFileInfo(c.Request.Context(), storage.BucketTypeTracks, objectPath)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to get file info: %v", err))
		return
	}

	if fileInfo.Size != request.ActualSize {
		jsonResponse.ResponseBadRequest(c, "File size mismatch")
		return
	}

	// Create song record
	baseModelInstance, _ := baseModel.NewBaseModel(h.generator)
	songData := &model.Song{
		BaseModel:        *baseModelInstance,
		ArtistID:         uploadSession.ArtistID,
		Title:            request.Title,
		Description:      request.Description,
		GenreID:          request.GenreID,
		MoodID:           request.MoodID,
		FileURL:          request.FileURL,
		FileSizeBytes:    &fileInfo.Size,
		DurationSeconds:  nil, // Will be set after processing
		ArtworkURL:       "",  // Can be added later
		Tier:             model.ContentTierPublicDiscovery,
		IsProcessed:      false,
		ProcessingStatus: model.ProcessingStatusPending,
		ProcessingError:  "",
		PlayCount:        0,
		LikeCount:        0,
		TipCount:         0,
		TotalTips:        0.0,
		IsActive:         true,
	}

	songID, err := h.musicService.CreateSong(c.Request.Context(), songData)
	if h.HandleError(c, err) {
		return
	}

	// Prepare metadata for Celery task
	metadata := queue.AudioProcessingMetadata{
		OriginalFilename: uploadSession.Filename,
		FileSize:         fileInfo.Size,
		ContentType:      h.getContentTypeFromFilename(uploadSession.Filename),
		UploadSessionID:  request.UploadID,
		Title:            request.Title,
		GenreID:          request.GenreID,
		MoodID:           request.MoodID,
		Description:      request.Description,
		AdditionalData: map[string]string{
			"user_agent": c.GetHeader("User-Agent"),
			"client_ip":  c.ClientIP(),
		},
	}

	// Create processing task
	callbackURL := fmt.Sprintf("/api/v1/processing/callback/%d", songID)
	processingTask := h.celeryClient.CreateProcessingTaskForSong(
		songID,
		uploadSession.ArtistID,
		objectPath,
		metadata,
		callbackURL,
	)

	// Apply custom processing config if provided
	if request.ProcessingConfig != nil {
		if request.ProcessingConfig.TargetLUFS != nil {
			processingTask.ProcessingConfig.TargetLUFS = *request.ProcessingConfig.TargetLUFS
		}
		if len(request.ProcessingConfig.GenerateFormats) > 0 {
			processingTask.ProcessingConfig.GenerateFormats = request.ProcessingConfig.GenerateFormats
		}
		if request.ProcessingConfig.QualityEnhancement != nil {
			processingTask.ProcessingConfig.QualityEnhancement = *request.ProcessingConfig.QualityEnhancement
		}
		if request.ProcessingConfig.PreserveDynamicRange != nil {
			processingTask.ProcessingConfig.PreserveDynamicRange = *request.ProcessingConfig.PreserveDynamicRange
		}
		if request.ProcessingConfig.ProcessingIntensity != "" {
			processingTask.ProcessingConfig.ProcessingIntensity = request.ProcessingConfig.ProcessingIntensity
		}
	}

	// Submit to Celery
	taskID, err := h.celeryClient.SubmitAudioProcessingTask(c.Request.Context(), processingTask)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to submit processing task: %v", err))
		return
	}

	// Update upload session status
	err = h.musicService.UpdateUploadSession(c.Request.Context(), request.UploadID, "completed")
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update upload session status: %v\n", err)
	}

	response := &CompleteUploadResponse{
		SongID:              songID,
		Status:              "upload_completed",
		Message:             "Upload completed successfully, processing started",
		ProcessingTaskID:    taskID,
		EstimatedCompletion: time.Now().Add(10 * time.Minute), // Estimated processing time
		UploadSummary: &UploadSummary{
			FileSize:   fileInfo.Size,
			Format:     h.getFormatFromFilename(uploadSession.Filename),
			UploadedAt: time.Now(),
		},
		NextSteps: []string{
			"Track is being processed for optimal quality",
			"You'll receive a notification when processing is complete",
			"Check processing status using the tracking URL",
			"Share your music once processing is done",
		},
		TrackingURL: fmt.Sprintf("/api/v1/processing/status/%s", taskID),
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *MusicHandler) GetUploadStatus(c *gin.Context) {
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		jsonResponse.ResponseBadRequest(c, "Upload ID is required")
		return
	}

	uploadSession, err := h.musicService.GetUploadSession(c.Request.Context(), uploadID)
	if h.HandleError(c, err) {
		return
	}

	if uploadSession == nil {
		jsonResponse.ResponseNotFound(c)
		return
	}

	response := map[string]interface{}{
		"upload_id":  uploadSession.ID,
		"status":     uploadSession.Status,
		"expires_at": uploadSession.ExpiresAt,
		"filename":   uploadSession.Filename,
		"file_size":  uploadSession.FileSize,
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *MusicHandler) GetProcessingStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		jsonResponse.ResponseBadRequest(c, "Task ID is required")
		return
	}

	// Get task result from Celery
	celeryResult, err := h.celeryClient.GetTaskResult(c.Request.Context(), taskID)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to get task result: %v", err))
		return
	}

	response := &ProcessingStatusResponse{
		TaskID:    taskID,
		Status:    celeryResult.Status,
		CreatedAt: time.Now(), // You might want to store this in Redis when creating the task
		UpdatedAt: time.Now(),
	}

	// Handle different status cases
	switch celeryResult.Status {
	case "PENDING":
		response.Progress = &ProcessingProgress{
			Stage:       "queued",
			Percentage:  0.0,
			CurrentStep: "Waiting in processing queue",
		}

	case "STARTED":
		response.Progress = &ProcessingProgress{
			Stage:       "processing",
			Percentage:  25.0,
			CurrentStep: "Audio processing in progress",
			ETA:         timePtr(time.Now().Add(8 * time.Minute)),
		}

	case "SUCCESS":
		// Parse the audio processing result
		audioResult, err := h.celeryClient.GetAudioProcessingResult(c.Request.Context(), taskID)
		if err != nil {
			jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to parse processing result: %v", err))
			return
		}
		
		response.Result = audioResult
		response.Progress = &ProcessingProgress{
			Stage:       "completed",
			Percentage:  100.0,
			CurrentStep: "Processing completed successfully",
		}

		// TODO: Update song record with processing results
		// h.updateSongWithProcessingResults(audioResult)

	case "FAILURE":
		response.Error = "Processing failed"
		if celeryResult.Traceback != "" {
			response.Error = celeryResult.Traceback
		}
		response.Progress = &ProcessingProgress{
			Stage:       "failed",
			Percentage:  0.0,
			CurrentStep: "Processing failed",
		}

	case "RETRY":
		response.Progress = &ProcessingProgress{
			Stage:       "retrying",
			Percentage:  10.0,
			CurrentStep: "Retrying processing due to temporary error",
		}
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *MusicHandler) GetStreamingURL(c *gin.Context) {
	songID := c.Param("song_id")
	if songID == "" {
		jsonResponse.ResponseBadRequest(c, "Song ID is required")
		return
	}
	format := c.DefaultQuery("format", "mp3_320")

	userTierInterface, exists := c.Get("user_tier")
	if !exists {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	userTier, ok := userTierInterface.(string)
	if !ok {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("invalid user tier type"))
		return
	}

	if !h.canAccessFormat(userTier, format) {
		jsonResponse.ResponseForbidden(c)
		return
	}

	objectPath := fmt.Sprintf("processed/%s/%s.%s", songID, format, h.getFormatExtension(format))

	streamingURL, err := h.storageService.GetStreamingURL(
		c.Request.Context(),
		objectPath,
		24*time.Hour,
	)
	if h.HandleError(c, err) {
		return
	}

	response := map[string]interface{}{
		"streaming_url": streamingURL,
		"format":        format,
		"expires_at":    time.Now().Add(24 * time.Hour),
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *MusicHandler) ProcessingCallback(c *gin.Context) {
	songID := c.Param("song_id")
	if songID == "" {
		jsonResponse.ResponseBadRequest(c, "Song ID is required")
		return
	}

	var callbackData queue.AudioProcessingResult
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		jsonResponse.ResponseBadRequest(c, "Invalid callback data: "+err.Error())
		return
	}

	// TODO: Update song record with processing results
	// This would include:
	// - Update processing_status to "completed" or "failed"
	// - Store audio analysis results
	// - Update duration_seconds
	// - Set is_processed = true
	// - Store processed file URLs

	fmt.Printf("Processing callback received for song %s: %+v\n", songID, callbackData)

	jsonResponse.ResponseOK(c, map[string]string{"status": "callback_received"})
}

// Helper methods

func (h *MusicHandler) canAccessFormat(userTier, format string) bool {
	switch userTier {
	case "free":
		return format == "mp3_320"
	case "premium":
		return format == "mp3_320" || format == "flac_cd"
	case "audiophile":
		return true // Can access all formats
	default:
		return false
	}
}

func (h *MusicHandler) isValidAudioFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validFormats := []string{".flac", ".wav", ".aiff", ".mp3"}

	for _, format := range validFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func (h *MusicHandler) generateUploadID(artistID uint64, filename string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("upload_%d_%d_%s", artistID, timestamp, h.sanitizeFilename(filename))
}

func (h *MusicHandler) sanitizeFilename(filename string) string {
	// Remove extension and special characters, keep only alphanumeric and hyphens
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name)
}

func (h *MusicHandler) extractObjectPathFromURL(fileURL string) string {
	// Extract object path from MinIO URL
	// This is a simplified implementation
	parts := strings.Split(fileURL, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return ""
}

func (h *MusicHandler) getFormatExtension(format string) string {
	switch format {
	case "mp3_320", "mp3_256", "mp3_192":
		return "mp3"
	case "flac_cd", "flac_hires":
		return "flac"
	default:
		return "mp3"
	}
}

func (h *MusicHandler) getFormatFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp3":
		return "MP3"
	case ".flac":
		return "FLAC"
	case ".wav":
		return "WAV"
	case ".aiff":
		return "AIFF"
	default:
		return "Unknown"
	}
}

func (h *MusicHandler) getContentTypeFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".flac":
		return "audio/flac"
	case ".wav":
		return "audio/wav"
	case ".aiff":
		return "audio/aiff"
	default:
		return "application/octet-stream"
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}