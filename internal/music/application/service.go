package application

import (
	"music-app-backend/internal/music/adapters/repository"
	model "music-app-backend/internal/music/domain"
	baseModel "music-app-backend/pkg/model"

	goflakeid "github.com/capy-engineer/go-flakeid"
)

type IMusicService interface {
	InsertArtist(artist *model.CreateArtistDTO) error
}

type MusicService struct {
	repository repository.IMusicRepository
	generator  *goflakeid.Generator
}

func NewMusicService(repository repository.IMusicRepository, generator *goflakeid.Generator) *MusicService {
	return &MusicService{repository: repository, generator: generator}
}

func (s *MusicService) InsertArtist(artist *model.CreateArtistDTO) error {
	_base, err := baseModel.NewBaseModel(s.generator)
	if err != nil {
		return err
	}
	artistModel := &model.Artist{
		BaseModel:       *_base,
		UserID:          artist.UserID,
		ArtistName:      artist.ArtistName,
		Bio:             artist.Bio,
		ProfileImageURL: artist.ProfileImageURL,
		BannerImageURL:  artist.BannerImageURL,
		WebsiteURL:      artist.WebsiteURL,
		SpotifyURL:      artist.SpotifyURL,
		InstagramURL:    artist.InstagramURL,
		TwitterURL:      artist.TwitterURL,
		YoutubeURL:      artist.YoutubeURL,
	}

	return s.repository.InsertArtist(artistModel)
}

// func (s *MusicService) InitiateUpload(ctx context.Context, request *model.InitiateUploadRequest) (*model.UploadResponse, error) {
// 	// Implementation for initiating an upload
// 	// This is a placeholder and should be replaced with actual logic
// 	// return &model.UploadResponse{
// 	// 	JobID: "example-job-id",
// 	// }, nil

// 	return nil, nil // Placeholder return
// }

// func (s *MusicService) CompleteUpload(ctx context.Context, jobID string, fileURL string) error {
// 	// Implementation for completing an upload
// 	// This is a placeholder and should be replaced with actual logic
// 	// return nil

// 	return nil // Placeholder return
// }

// func (s *MusicService) GetUploadStatus(ctx context.Context, jobID string) (*model.UploadStatusResponse, error) {
// 	// Implementation for getting the status of an upload
// 	// This is a placeholder and should be replaced with actual logic
// 	// return &model.UploadStatusResponse{
// 	// 	Status: "completed",
// 	// }, nil

// 	return nil, nil // Placeholder return
// }
