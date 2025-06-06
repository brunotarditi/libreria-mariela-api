package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(65);not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"-"`
}
