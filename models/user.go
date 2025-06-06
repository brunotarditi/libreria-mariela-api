package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Email     string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	LastLogin time.Time `json:"last_login"`
	Roles     []Role    `gorm:"many2many:user_roles;"`
}
