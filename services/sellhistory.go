package services

import (
	"fmt"
	"libreria/constants"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"
	"libreria/utils"

	"gorm.io/gorm"
)

type SellHistoryService interface {
	CreateSell(request requests.SellHistoryRequest) (models.SellHistory, error)
	DeleteSell(id string) error
}

type sellHistoryService struct {
	db                   *gorm.DB
	sellRepo             repositories.SellHistoryRepository
	productStockRepo     repositories.ProductStockRepository
	stockMovementRepo    repositories.StockMovementRepository
	productStockService  ProductStockService
	stockMovementService StockMovementService
}

func NewSellHistoryService(db *gorm.DB, sellRepo repositories.SellHistoryRepository, productStockRepo repositories.ProductStockRepository, stockMovementRepo repositories.StockMovementRepository, productStockService ProductStockService, stockMovementService StockMovementService) SellHistoryService {
	return &sellHistoryService{
		db:                   db,
		sellRepo:             sellRepo,
		productStockRepo:     productStockRepo,
		stockMovementRepo:    stockMovementRepo,
		productStockService:  productStockService,
		stockMovementService: stockMovementService,
	}
}

func (s *sellHistoryService) CreateSell(request requests.SellHistoryRequest) (models.SellHistory, error) {
	if err := request.Validate(s.db); err != nil {
		return models.SellHistory{}, err
	}

	averageCost, stock, err := utils.CalculateAverageCostAndStock(s.db, request.ProductID)

	if err != nil {
		return models.SellHistory{}, err
	}

	if stock < int64(request.Quantity) {
		return models.SellHistory{}, fmt.Errorf("stock insuficiente para producto %d", request.ProductID)
	}

	var product models.Product
	if err := s.db.First(&product, request.ProductID).Error; err != nil {
		return models.SellHistory{}, err
	}

	sell, err := request.ToModel()
	if err != nil {
		return models.SellHistory{}, err
	}

	sell.AverageCost = averageCost
	sell.Price = averageCost * (1 + product.ProfitMargin/100)

	tx := s.db.Begin()
	if err := s.sellRepo.Create(&sell); err != nil {
		tx.Rollback()
		return models.SellHistory{}, err
	}

	if err := applyMovementFlow(s.productStockService, s.stockMovementService, request.ProductID, request.Quantity, constants.STOCK_MOVEMENT_TYPE_OUT, sell.ID, "Nueva venta"); err != nil {
		tx.Rollback()
		return models.SellHistory{}, err
	}

	tx.Commit()
	return sell, nil
}

func (s *sellHistoryService) DeleteSell(id string) error {
	tx := s.db.Begin()
	sell, err := s.sellRepo.FindByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := applyMovementFlow(s.productStockService, s.stockMovementService, sell.ProductID, sell.Quantity, constants.STOCK_MOVEMENT_TYPE_IN, sell.ID, "Devolución de venta"); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.sellRepo.Delete(id); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
