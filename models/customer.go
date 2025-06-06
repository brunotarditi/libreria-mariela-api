package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name          string        `gorm:"type:varchar(65);not null" json:"name"`
	ContactInfo   string        `gorm:"type:varchar(255)" json:"contact_info"`
	SellHistories []SellHistory `gorm:"foreignKey:CustomerID" json:"-"`
}
