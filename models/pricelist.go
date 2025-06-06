package models

import (
	"time"

	"gorm.io/gorm"
)

type PriceList struct {
	gorm.Model
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	EffectiveAt time.Time `gorm:"not null" json:"effective_at"`
	IsActive    bool      `gorm:"default:true" json:"is_active"` // permite mantener historial
}
