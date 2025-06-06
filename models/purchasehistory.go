package models

import (
	"gorm.io/gorm"
)

type PurchaseHistory struct {
	gorm.Model
	ProductID  uint     `gorm:"not null" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product"`
	SupplierID uint     `gorm:"not null" json:"supplier_id"`
	Supplier   Supplier `gorm:"foreignKey:SupplierID" json:"supplier"`
	Cost       float64  `gorm:"type:decimal(10,2);not null" json:"cost"`
	Quantity   int      `gorm:"not null" json:"quantity"`
}
