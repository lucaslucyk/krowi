package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDBaseModel struct {
	ID        uuid.UUID  `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (base *UUIDBaseModel) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("ID", uuid.String())
	return nil
}
