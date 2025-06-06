package models

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(65);not null" json:"name"`
	Products []Product `gorm:"foreignKey:BrandID" json:"-"`
}
