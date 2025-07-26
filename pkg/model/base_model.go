package model

import (
	"time"

	goflakeid "github.com/hxuan190/go-flakeid"
)

const (
	EntityUser             goflakeid.EntityType = 0
	EntityArtist           goflakeid.EntityType = 1
	EntityGenre            goflakeid.EntityType = 2
	EntityMood             goflakeid.EntityType = 3
	EntitySong             goflakeid.EntityType = 4
	EntitySongPlay         goflakeid.EntityType = 5
	EntityUserFavorite     goflakeid.EntityType = 6
	EntityArtistFollower   goflakeid.EntityType = 7
	EntityPlaylist         goflakeid.EntityType = 8
	EntityPlaylistSong     goflakeid.EntityType = 9
	EntityTip              goflakeid.EntityType = 10
	EntityArtistMessage    goflakeid.EntityType = 11
	EntityMessageDelivery  goflakeid.EntityType = 12
	EntityListeningSession goflakeid.EntityType = 13
	EntityUserPreference   goflakeid.EntityType = 14
	EntityDailyArtistStat  goflakeid.EntityType = 15
	EntityDailySongStat    goflakeid.EntityType = 16
)

type BaseModel struct {
	ID        uint64 `json:"id" gorm:"primaryKey;"`
	Status    string `json:"status" gorm:"default:active;check:status IN ('active', 'inactive')"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewBaseModel(generator *goflakeid.IDGenerator, entiry goflakeid.EntityType) (*BaseModel, error) {
	id, err := generator.Generate(entiry)
	if err != nil {
		return nil, err
	}

	return &BaseModel{
		ID:        id,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}, nil
}
