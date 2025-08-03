package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	Email      string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	FirstName  string    `gorm:"type:varchar(50)" json:"first_name"`
	LastName   string    `gorm:"type:varchar(50)" json:"last_name"`
	BirthDate  time.Time `json:"birth_date"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	LastLogin  time.Time `json:"last_login"`
	Roles      []Role    `gorm:"many2many:user_roles;"`
}
