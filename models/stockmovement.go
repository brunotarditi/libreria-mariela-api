package models

import "gorm.io/gorm"

type StockMovement struct {
	gorm.Model
	ProductID    uint    `gorm:"not null" json:"product_id"`
	Product      Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity     int     `gorm:"not null" json:"quantity"`
	MovementType int     `gorm:"not null" json:"movement_type"` // "purchase", "sale", "adjustment"
	ReferenceID  *uint   `json:"reference_id"`                  // opcional: puede ser ID de venta o compra
	Note         string  `gorm:"type:varchar(150)" json:"note"`
}
