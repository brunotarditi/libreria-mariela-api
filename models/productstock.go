package models

import "gorm.io/gorm"

type ProductStock struct {
	gorm.Model
	ProductID uint    `gorm:"not null;unique" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null" json:"quantity"` // stock actual
}
