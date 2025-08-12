// pkg/storage/minio.go
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	BucketTypeTracks    = "audora-tracks"
	BucketTypeProcessed = "processed-tracks"
	BucketTypeGeneral   = "audora"
)

type MinIOConfig struct {
	Endpoint          string
	AccessKey         string
	SecretKey         string
	UseSSL            bool
	BucketName        string
	TracksBucket      string
	ProcessedBucket   string
	PipelineAccessKey string
	PipelineSecretKey string
}

type MinIOService struct {
	client         *minio.Client
	pipelineClient *minio.Client
	config         *MinIOConfig
}

type UploadResult struct {
	BucketName string    `json:"bucket_name"`
	ObjectName string    `json:"object_name"`
	Size       int64     `json:"size"`
	ETag       string    `json:"etag"`
	UploadedAt time.Time `json:"uploaded_at"`
	URL        string    `json:"url"`
}

type PresignedUploadURL struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt time.Time         `json:"expires_at"`
}

func NewMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:          getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey:         getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey:         getEnv("MINIO_SECRET_KEY", "minioadmin"),
		UseSSL:            getEnv("MINIO_USE_SSL", "false") == "true",
		BucketName:        getEnv("MINIO_BUCKET_NAME", "audora"),
		TracksBucket:      getEnv("MINIO_TRACKS_BUCKET", "audora-tracks"),
		ProcessedBucket:   getEnv("MINIO_PROCESSED_BUCKET", "processed-tracks"),
		PipelineAccessKey: getEnv("MINIO_PIPELINE_ACCESS_KEY", "pipeline-user"),
		PipelineSecretKey: getEnv("MINIO_PIPELINE_SECRET_KEY", "pipeline-secret-key"),
	}
}

func NewMinIOService(config *MinIOConfig) (*MinIOService, error) {
	// Main client for general operations
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Pipeline client for audio processing operations
	pipelineClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.PipelineAccessKey, config.PipelineSecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO pipeline client: %w", err)
	}

	service := &MinIOService{
		client:         client,
		pipelineClient: pipelineClient,
		config:         config,
	}

	// Ensure buckets exist
	if err := service.ensureBucketsExist(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure buckets exist: %w", err)
	}

	return service, nil
}

func (s *MinIOService) ensureBucketsExist(ctx context.Context) error {
	buckets := []string{
		s.config.BucketName,
		s.config.TracksBucket,
		s.config.ProcessedBucket,
	}

	for _, bucket := range buckets {
		exists, err := s.client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("failed to check if bucket %s exists: %w", bucket, err)
		}

		if !exists {
			err = s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
	}

	return nil
}

// GetPresignedUploadURL generates a presigned URL for direct client uploads
func (s *MinIOService) GetPresignedUploadURL(ctx context.Context, objectName string, bucketType string, expiry time.Duration) (*PresignedUploadURL, error) {
	bucket := s.getBucketName(bucketType)

	// Generate presigned PUT URL
	url, err := s.client.PresignedPutObject(ctx, bucket, objectName, expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &PresignedUploadURL{
		URL:    url.String(),
		Method: "PUT",
		Headers: map[string]string{
			"Content-Type": "application/octet-stream",
		},
		ExpiresAt: time.Now().Add(expiry),
	}, nil
}

// UploadFile uploads a file to MinIO
func (s *MinIOService) UploadFile(ctx context.Context, bucketType, objectName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	bucket := s.getBucketName(bucketType)

	info, err := s.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate access URL
	url := s.getObjectURL(bucket, objectName)

	return &UploadResult{
		BucketName: bucket,
		ObjectName: objectName,
		Size:       info.Size,
		ETag:       info.ETag,
		UploadedAt: time.Now(),
		URL:        url,
	}, nil
}

// DownloadFile downloads a file from MinIO
func (s *MinIOService) DownloadFile(ctx context.Context, bucketType, objectName string) (io.ReadCloser, error) {
	bucket := s.getBucketName(bucketType)

	object, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return object, nil
}

// PipelineUploadFile uploads processed files using pipeline credentials
func (s *MinIOService) PipelineUploadFile(ctx context.Context, bucketType, objectName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	bucket := s.getBucketName(bucketType)

	info, err := s.pipelineClient.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload processed file: %w", err)
	}

	url := s.getObjectURL(bucket, objectName)

	return &UploadResult{
		BucketName: bucket,
		ObjectName: objectName,
		Size:       info.Size,
		ETag:       info.ETag,
		UploadedAt: time.Now(),
		URL:        url,
	}, nil
}

// GetStreamingURL generates a presigned URL for streaming
func (s *MinIOService) GetStreamingURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.config.ProcessedBucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate streaming URL: %w", err)
	}

	return url.String(), nil
}

// DeleteFile deletes a file from MinIO
func (s *MinIOService) DeleteFile(ctx context.Context, bucketType, objectName string) error {
	bucket := s.getBucketName(bucketType)

	err := s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListFiles lists files in a bucket with a prefix
func (s *MinIOService) ListFiles(ctx context.Context, bucketType, prefix string) ([]minio.ObjectInfo, error) {
	bucket := s.getBucketName(bucketType)

	objectCh := s.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// GetFileInfo gets information about a file
func (s *MinIOService) GetFileInfo(ctx context.Context, bucketType, objectName string) (*minio.ObjectInfo, error) {
	bucket := s.getBucketName(bucketType)

	info, err := s.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &info, nil
}

// CopyFile copies a file from one location to another
func (s *MinIOService) CopyFile(ctx context.Context, srcBucketType, srcObjectName, destBucketType, destObjectName string) error {
	srcBucket := s.getBucketName(srcBucketType)
	destBucket := s.getBucketName(destBucketType)

	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObjectName,
	}

	dest := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: destObjectName,
	}

	_, err := s.pipelineClient.CopyObject(ctx, dest, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// GenerateUploadPath generates a structured path for uploads
func (s *MinIOService) GenerateUploadPath(artistID uint64, filename string) string {
	// Clean filename and ensure proper extension
	cleanName := filepath.Base(filename)
	ext := strings.ToLower(filepath.Ext(cleanName))
	nameWithoutExt := strings.TrimSuffix(cleanName, ext)

	// Generate timestamp-based path to avoid conflicts
	timestamp := time.Now().Format("2006/01/02")

	return fmt.Sprintf("uploads/%d/%s/%s%s", artistID, timestamp, nameWithoutExt, ext)
}

// GenerateProcessedPath generates a structured path for processed files
func (s *MinIOService) GenerateProcessedPath(artistID uint64, songID uint64, format, quality string) string {
	timestamp := time.Now().Format("2006/01")
	return fmt.Sprintf("processed/%d/%s/%d/%s/%s", artistID, timestamp, songID, format, quality)
}

// Helper methods
func (s *MinIOService) getBucketName(bucketType string) string {
	switch bucketType {
	case "tracks", "original":
		return s.config.TracksBucket
	case "processed", "streaming":
		return s.config.ProcessedBucket
	case "general", "public":
		return s.config.BucketName
	default:
		return s.config.BucketName
	}
}

func (s *MinIOService) getObjectURL(bucket, objectName string) string {
	scheme := "http"
	if s.config.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.config.Endpoint, bucket, objectName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
