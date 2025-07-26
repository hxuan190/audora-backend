package model

import (
	"time"

	goflakeid "github.com/hxuan190/go-flakeid"
)

const (
	EntityUser            goflakeid.EntityType = 0
	EntityArtist          goflakeid.EntityType = 1
	EntitySong            goflakeid.EntityType = 2
	EntityPlaylist        goflakeid.EntityType = 3
	EntityTip             goflakeid.EntityType = 4
	EntityMessage         goflakeid.EntityType = 5
	EntityMessageDelivery goflakeid.EntityType = 6
	EntitySession         goflakeid.EntityType = 7
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
