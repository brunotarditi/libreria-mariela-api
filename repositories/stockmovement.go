package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(stock *models.StockMovement) error
	WithTx(tx *gorm.DB) StockMovementRepository
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) WithTx(tx *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: tx}
}

func (r *stockMovementRepository) Create(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}
