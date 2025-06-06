package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type SellHistoryRepository interface {
	Create(sell *models.SellHistory) error
	FindByID(id string) (models.SellHistory, error)
	Delete(id string) error
}

type sellHistoryRepository struct {
	db *gorm.DB
}

func NewSellHistoryRepository(db *gorm.DB) SellHistoryRepository {
	return &sellHistoryRepository{db: db}
}

func (r *sellHistoryRepository) Create(Sell *models.SellHistory) error {
	return r.db.Create(Sell).Error
}

func (r *sellHistoryRepository) FindByID(id string) (models.SellHistory, error) {
	var Sell models.SellHistory
	err := r.db.First(&Sell, id).Error
	return Sell, err
}

func (r *sellHistoryRepository) Delete(id string) error {
	return r.db.Delete(&models.SellHistory{}, id).Error
}
