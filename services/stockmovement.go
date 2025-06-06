package services

import (
	"libreria/models"
	"libreria/repositories"

	"gorm.io/gorm"
)

type StockMovementService interface {
	ApplyMovement(stockMovement models.StockMovement) error
}

type stockMovementService struct {
	db                *gorm.DB
	stockMovementRepo repositories.StockMovementRepository
}

func NewStockMovementService(db *gorm.DB, stockRepo repositories.StockMovementRepository) StockMovementService {
	return &stockMovementService{db: db, stockMovementRepo: stockRepo}
}

func (s *stockMovementService) ApplyMovement(stockMovement models.StockMovement) error {

	tx := s.db.Begin()

	if err := s.stockMovementRepo.Create(&stockMovement); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
