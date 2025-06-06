package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(stock *models.StockMovement) error
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}
