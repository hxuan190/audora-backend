package model

import (
	"time"

	goflakeid "github.com/capy-engineer/go-flakeid"
)

type BaseModel struct {
	ID        uint64 `json:"id" gorm:"primaryKey;"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewBaseModel(generator *goflakeid.Generator) (*BaseModel, error) {
	id, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	return &BaseModel{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
