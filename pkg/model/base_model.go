package model

import (
	"time"

	goflakeid "github.com/capy-engineer/go-flakeid"
)

type BaseModel struct {
	ID        uint64 `json:"id" gorm:"primaryKey;"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewBaseModel(generator *goflakeid.Generator) (*BaseModel, error) {
	id, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	return &BaseModel{
		ID:        id,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}, nil
}
