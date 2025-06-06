package services

import (
	"fmt"
	"libreria/constants"
	"libreria/models"
	"libreria/repositories"

	"gorm.io/gorm"
)

type ProductStockService interface {
	ApplyMovement(productStock models.ProductStock, movementType int) error
}

type productStockService struct {
	db               *gorm.DB
	productStockRepo repositories.ProductStockRepository
}

func NewProductStockService(db *gorm.DB, stockRepo repositories.ProductStockRepository) ProductStockService {
	return &productStockService{db: db, productStockRepo: stockRepo}
}

func (s *productStockService) ApplyMovement(productStock models.ProductStock, movementType int) error {

	tx := s.db.Begin()
	stockExist, err := s.productStockRepo.FindByID(productStock.ProductID)
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}
	if err == gorm.ErrRecordNotFound {

		if err := s.productStockRepo.Create(&productStock); err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	} else {
		if movementType == int(constants.STOCK_MOVEMENT_TYPE_IN) {
			stockExist.Quantity += productStock.Quantity
		}

		if movementType == int(constants.STOCK_MOVEMENT_TYPE_OUT) {
			stockExist.Quantity -= productStock.Quantity
		}

		if stockExist.Quantity < 0 {
			return fmt.Errorf("stock insuficiente para producto %d", productStock.ProductID)
		}

		if err := s.productStockRepo.Update(&stockExist); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
}
