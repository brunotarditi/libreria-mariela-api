package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `json:"user_id"` // null si no hay usuario
	Route     string    `gorm:"type:varchar(255);not null;index:idx_route_method"`
	Method    string    `gorm:"type:varchar(10);not null;index:idx_route_method"`
	IP        string    `gorm:"type:varchar(100);not null"`
	RequestAt time.Time `gorm:"not null;index:idx_request_at"`
}
