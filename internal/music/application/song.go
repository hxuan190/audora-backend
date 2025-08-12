package application

import (
	"context"
	model "music-app-backend/internal/music/domain"
	baseModel "music-app-backend/pkg/model"
	"music-app-backend/pkg/queue"
)

func (s *MusicService) CreateSong(ctx context.Context, song *model.Song) (uint64, error) {
	err := s.repository.InsertSong(song)
	if err != nil {
		return 0, err
	}

	return song.ID, nil
}

func (s *MusicService) UpdateSongWithProcessingResults(ctx context.Context, songID uint64, result *queue.AudioProcessingResult) error {
	// Prepare update data
	updates := make(map[string]interface{})

	if result.Success {
		updates["processing_status"] = model.ProcessingStatusCompleted
		updates["is_processed"] = true
		updates["processing_error"] = ""

		// Update duration if available from audio analysis
		if result.AudioAnalysis.Duration > 0 {
			durationSeconds := int(result.AudioAnalysis.Duration)
			updates["duration_seconds"] = durationSeconds
		}

		// Store audio analysis results as JSON strings or separate fields
		if result.AudioAnalysis.OriginalFormat != "" {
			updates["bpm"] = result.AudioAnalysis.SpectralCentroid // This could be BPM if available
		}

		// Store processed audio formats
		if len(result.ProcessedFormats) > 0 {
			err := s.storeProcessedAudioFormats(ctx, songID, result.ProcessedFormats)
			if err != nil {
				return err
			}
		}

		// Store audio analysis
		err := s.storeAudioAnalysis(ctx, songID, &result.AudioAnalysis, result.QualityScore, result.ProcessingTime)
		if err != nil {
			return err
		}

	} else {
		updates["processing_status"] = model.ProcessingStatusFailed
		updates["is_processed"] = false
		if result.Error != "" {
			updates["processing_error"] = result.Error
		}
	}

	return s.repository.UpdateSongProcessingResult(ctx, songID, updates)
}

func (s *MusicService) storeProcessedAudioFormats(ctx context.Context, songID uint64, formats []queue.ProcessedAudioFormat) error {
	if len(formats) == 0 {
		return nil
	}

	processedFormats := make([]model.ProcessedAudioFormat, len(formats))
	for i, format := range formats {
		baseModelInstance, err := s.generateBaseModel()
		if err != nil {
			return err
		}

		processedFormats[i] = model.ProcessedAudioFormat{
			BaseModel:    *baseModelInstance,
			SongID:       songID,
			Format:       format.Format,
			ObjectPath:   format.ObjectPath,
			FileSize:     format.FileSize,
			Bitrate:      &format.Bitrate,
			SampleRate:   &format.SampleRate,
			BitDepth:     &format.BitDepth,
			Duration:     format.Duration,
			QualityScore: format.QualityScore,
		}
	}

	return s.repository.CreateProcessedAudioFormats(ctx, processedFormats)
}

func (s *MusicService) storeAudioAnalysis(ctx context.Context, songID uint64, analysis *queue.AudioAnalysis, qualityScore, processingTime float64) error {
	baseModelInstance, err := s.generateBaseModel()
	if err != nil {
		return err
	}

	audioAnalysis := &model.AudioAnalysis{
		BaseModel:          *baseModelInstance,
		SongID:             songID,
		OriginalFormat:     analysis.OriginalFormat,
		OriginalBitrate:    &analysis.OriginalBitrate,
		OriginalSampleRate: &analysis.OriginalSampleRate,
		OriginalBitDepth:   &analysis.OriginalBitDepth,
		Duration:           analysis.Duration,
		OriginalLUFS:       &analysis.OriginalLUFS,
		ProcessedLUFS:      &analysis.ProcessedLUFS,
		DynamicRange:       &analysis.DynamicRange,
		PeakLevel:          &analysis.PeakLevel,
		TruePeak:           &analysis.TruePeak,
		SpectralCentroid:   &analysis.SpectralCentroid,
		THDPlusN:           &analysis.THDPlusN,
		StereoWidth:        &analysis.StereoWidth,
		HasClipping:        analysis.HasClipping,
		HasArtifacts:       analysis.HasArtifacts,
		QualityGrade:       analysis.QualityGrade,
		QualityScore:       qualityScore,
		ProcessingTime:     processingTime,
		Warnings:           []string{}, // Initialize empty slice for now
	}

	return s.repository.CreateAudioAnalysis(ctx, audioAnalysis)
}

func (s *MusicService) generateBaseModel() (*baseModel.BaseModel, error) {
	return baseModel.NewBaseModel(s.generator)
}

func (s *MusicService) GetSongByID(ctx context.Context, songID uint64) (*model.Song, error) {
	return s.repository.GetSongByID(ctx, songID)
}
