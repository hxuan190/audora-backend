package application

import (
	"context"
	model "music-app-backend/internal/music/domain"
)

func (s *MusicService) CreateSong(ctx context.Context, song *model.Song) (uint64, error) {
	err := s.repository.InsertSong(song)
	if err != nil {
		return 0, err
	}

	return song.ID, nil
}
