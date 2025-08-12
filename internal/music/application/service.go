package application

import (
	"context"
	"music-app-backend/internal/music/adapters/repository"
	model "music-app-backend/internal/music/domain"
	baseModel "music-app-backend/pkg/model"

	goflakeid "github.com/capy-engineer/go-flakeid"
)

type IMusicService interface {
	InsertArtist(artist *model.CreateArtistDTO) error
	CreateUploadSession(upload *model.UploadSession) error
	CreateSong(ctx context.Context, song *model.Song) (uint64, error)
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
