package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uint
	User      User
	TokenHash []byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
