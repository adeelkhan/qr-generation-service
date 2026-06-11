package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QRCode struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	InputText string    `gorm:"type:text;not null"`
	ImageData []byte    `gorm:"type:bytea;not null"`
	Title     string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
}

func (q *QRCode) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
