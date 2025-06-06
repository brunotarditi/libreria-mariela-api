package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `json:"user_id"` // null si no hay usuario
	Route     string    `gorm:"type:varchar(255);not null"`
	Method    string    `gorm:"type:varchar(10);not null"`
	IP        string    `gorm:"type:varchar(100);not null"`
	RequestAt time.Time `gorm:"not null"`
}
