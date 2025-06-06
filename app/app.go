package app

import (
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func NewApp(db *gorm.DB) *App {
	return &App{DB: db}
}
