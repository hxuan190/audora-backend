package http

import (
	"fmt"
	"music-app-backend/internal/music/application"
	model "music-app-backend/internal/music/domain"
	appError "music-app-backend/pkg/error"
	jsonResponse "music-app-backend/pkg/json"
	"music-app-backend/pkg/storage"
	"path/filepath"
	"strings"
	"time"

	baseModel "music-app-backend/pkg/model"

	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	musicService   *application.MusicService
	storageService *storage.MinIOService
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

	// Song metadata
	Description string `json:"description"`
}

type ProcessingIntensity string

const (
	ProcessingConservative ProcessingIntensity = "conservative" // Gentle processing, preserve dynamics
	ProcessingStandard     ProcessingIntensity = "standard"     // Balanced optimization
	ProcessingAggressive   ProcessingIntensity = "aggressive"   // Maximum loudness/presence
)

type CollaboratorRequest struct {
	ArtistID     *uint64  `json:"artist_id,omitempty"`     // If they're on the platform
	Name         string   `json:"name" binding:"required"` // Display name
	Role         string   `json:"role" binding:"required"` // "vocalist", "producer", "writer", etc.
	Email        *string  `json:"email,omitempty"`         // For future invitation
	SplitPercent *float64 `json:"split_percent,omitempty"` // Revenue split percentage
}

type CreditRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"` // "Mixed by", "Mastered by", "Recorded at", etc.
}

type ProcessingJob struct {
	SongID           uint64 `json:"song_id"`
	ArtistID         uint64 `json:"artist_id"`
	SourceObjectPath string `json:"source_object_path"`
}

type CompleteUploadResponse struct {
	SongID              uint64         `json:"song_id"`
	Status              string         `json:"status"`
	Message             string         `json:"message"`
	ProcessingJobID     string         `json:"processing_job_id"`
	EstimatedCompletion time.Time      `json:"estimated_completion"`
	UploadSummary       *UploadSummary `json:"upload_summary"`
	NextSteps           []string       `json:"next_steps"`
}

type UploadSummary struct {
	FileSize   int64     `json:"file_size"`
	Format     string    `json:"format"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewUploadHandler(musicService *application.MusicService, storageService *storage.MinIOService, generator *goflakeid.Generator) *UploadHandler {
	return &UploadHandler{
		musicService:   musicService,
		storageService: storageService,
		generator:      generator,
	}
}

func (h *UploadHandler) HandleError(c *gin.Context, err error) bool {
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

func (h *UploadHandler) InitiateUpload(c *gin.Context) {
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

	if !h.isValidAudioFormat(request.Filename) {
		jsonResponse.ResponseBadRequest(c, "Unsupported file format. Supported: FLAC, WAV, AIFF, MP3")
		return
	}

	maxSize := int64(600 * 1024 * 1024) // 600MB

	if request.FileSize <= 0 || request.FileSize > maxSize {
		jsonResponse.ResponseBadRequest(c, "File size must be between 1 byte and 600MB")
		return
	}

	uploadID := h.generateUploadID(request.ArtistID, request.Filename)
	objectPath := h.storageService.GenerateUploadPath(request.ArtistID, request.Filename)

	presignedURL, err := h.storageService.GetPresignedUploadURL(c.Request.Context(), objectPath, storage.BucketTypeTracks, 15*time.Minute)
	if err != nil {
		jsonResponse.ResponseInternalError(c, fmt.Errorf("failed to generate presigned URL: %v", err))
		return
	}

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

func (h *UploadHandler) CompleteUpload(c *gin.Context) {
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

	uploadSession, err := h.musicService.GetUploadSession(c.Request.Context(), request.UploadID)
	if h.HandleError(c, err) {
		return
	}

	if uploadSession.UserID != userID.(uint64) {
		jsonResponse.ResponseForbidden(c)
		return
	}

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

	baseModelInstance, _ := baseModel.NewBaseModel(h.generator)
	songData := &model.Song{
		BaseModel:        *baseModelInstance,
		ArtistID:         userID.(uint64),
		Title:            request.Title,
		Description:      request.Description,
		GenreID:          request.GenreID,
		MoodID:           request.MoodID,
		FileURL:          request.FileURL,
		FileSizeBytes:    &fileInfo.Size,
		DurationSeconds:  nil,                              // Duration will be set later after processing
		ArtworkURL:       "",                               // Artwork handling can be added later
		Tier:             model.ContentTierPublicDiscovery, // Default tier, can be overridden later
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

	// processingJob := &ProcessingJob{
	// 	SongID:           songID,
	// 	ArtistID:         userID.(uint64),
	// 	SourceObjectPath: objectPath,
	// }

	// TODO: QUEUE PROCESSING JOB
	// This is where you would enqueue the processing job to a message queue or worker system
	// For now, we will just simulate it
	processingJobID := fmt.Sprintf("processing_%d_%d", songID, time.Now().Unix())
	fmt.Println("Enqueued processing job:", processingJobID)

	err = h.musicService.UpdateUploadSession(c.Request.Context(), request.UploadID, "completed")
	if h.HandleError(c, err) {
		return
	}

	response := &CompleteUploadResponse{
		SongID:              songID,
		Status:              "completed",
		Message:             "Upload completed successfully",
		ProcessingJobID:     processingJobID,
		EstimatedCompletion: time.Now().Add(30 * time.Minute), // Simulated estimated completion
		UploadSummary: &UploadSummary{
			FileSize:   fileInfo.Size,
			Format:     h.getFormatExtension("mp3_320"),
			UploadedAt: time.Now(),
		},
		NextSteps: []string{
			"Share your song on social media",
			"Add collaborators or credits",
			"Set up monetization options",
			"Submit to playlists or radio stations",
			"Engage with your audience",
		},
	}
	jsonResponse.ResponseOK(c, response)
}

func (h *UploadHandler) GetUploadStatus(c *gin.Context) {
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
	}

	jsonResponse.ResponseOK(c, response)
}

func (h *UploadHandler) GetStreamingURL(c *gin.Context) {
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

// Helper function to validate audio file formats
func (h *UploadHandler) canAccessFormat(userTier, format string) bool {
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

func (h *UploadHandler) isValidAudioFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validFormats := []string{".flac", ".wav", ".aiff", ".mp3"}

	for _, format := range validFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func (h *UploadHandler) generateUploadID(artistID uint64, filename string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("upload_%d_%d_%s", artistID, timestamp, h.sanitizeFilename(filename))
}

func (h *UploadHandler) sanitizeFilename(filename string) string {
	// Remove extension and special characters, keep only alphanumeric and hyphens
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name)
}

func (h *UploadHandler) extractObjectPathFromURL(fileURL string) string {
	// Extract object path from MinIO URL
	// This is a simplified implementation
	parts := strings.Split(fileURL, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return ""
}

func (h *UploadHandler) getFormatExtension(format string) string {
	switch format {
	case "mp3_320", "mp3_256", "mp3_192":
		return "mp3"
	case "flac_cd", "flac_hires":
		return "flac"
	default:
		return "mp3"
	}
}
