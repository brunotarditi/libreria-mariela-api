package app

import (
	"github.com/brunotarditi/peak-auth/sdk/go"
	"gorm.io/gorm"
)

type App struct {
	DB             *gorm.DB
	PeakAuthClient *peakauth.Client
}

func NewApp(db *gorm.DB, peakAuthClient *peakauth.Client) *App {
	return &App{
		DB:             db,
		PeakAuthClient: peakAuthClient,
	}
}
