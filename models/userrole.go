package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"not null" json:"user_id"`
	RoleID uint `gorm:"not null" json:"role_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
	Role   Role `gorm:"foreignKey:RoleID" json:"role"`
}
