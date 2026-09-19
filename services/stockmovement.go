package services

import (
	"libreria/models"
	"libreria/repositories"

	"gorm.io/gorm"
)

type StockMovementService interface {
	ApplyMovement(stockMovement models.StockMovement) error
	WithTx(tx *gorm.DB) StockMovementService
}

type stockMovementService struct {
	db                *gorm.DB
	stockMovementRepo repositories.StockMovementRepository
}

func NewStockMovementService(db *gorm.DB, stockRepo repositories.StockMovementRepository) StockMovementService {
	return &stockMovementService{db: db, stockMovementRepo: stockRepo}
}

func (s *stockMovementService) WithTx(tx *gorm.DB) StockMovementService {
	return &stockMovementService{
		db:                tx,
		stockMovementRepo: s.stockMovementRepo.WithTx(tx),
	}
}

func (s *stockMovementService) ApplyMovement(stockMovement models.StockMovement) error {
	return s.stockMovementRepo.Create(&stockMovement)
}
