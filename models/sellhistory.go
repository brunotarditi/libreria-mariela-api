package models

import "gorm.io/gorm"

type SellHistory struct {
	gorm.Model
	ProductID   uint     `gorm:"not null" json:"product_id"`
	Product     Product  `gorm:"foreignKey:ProductID" json:"product"`
	CustomerID  uint     `gorm:"not null" json:"customer_id"`
	Customer    Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	Price       float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity    int      `gorm:"not null" json:"quantity"`
	AverageCost float64  `gorm:"type:decimal(10,2);not null" json:"average_cost"`
}
