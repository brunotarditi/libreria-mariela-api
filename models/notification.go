package models

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"` // null para notificaciones globales del sistema
	Title     string    `gorm:"type:varchar(120);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Type      string    `gorm:"type:varchar(20);not null;default:'info'" json:"type"` // 'info', 'warning', 'success', 'alert'
	Icon      string    `gorm:"type:varchar(50);not null;default:'notifications'" json:"icon"`
	Route     *string   `gorm:"type:varchar(150)" json:"route"` // ruta opcional ej: '/products'
	IsRead    bool      `gorm:"not null;default:false;index" json:"is_read"`
	CreatedAt time.Time `gorm:"not null;index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
