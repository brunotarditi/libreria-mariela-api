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
	DeleteSell(id uint64) error
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

	var sell models.SellHistory
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txSellRepo := s.sellRepo.WithTx(tx)
		txProductStockService := s.productStockService.WithTx(tx)
		txStockMovementService := s.stockMovementService.WithTx(tx)

		averageCost, stock, err := utils.CalculateAverageCostAndStock(tx, request.ProductID)
		if err != nil {
			return err
		}

		if stock < int64(request.Quantity) {
			return fmt.Errorf("stock insuficiente para producto %d", request.ProductID)
		}

		var product models.Product
		if err := tx.First(&product, request.ProductID).Error; err != nil {
			return err
		}

		model, err := request.ToModel()
		if err != nil {
			return err
		}
		sell = model
		sell.AverageCost = averageCost
		sell.Price = averageCost * (1 + product.ProfitMargin/100)

		if err := txSellRepo.Create(&sell); err != nil {
			return err
		}

		return applyMovementFlow(txProductStockService, txStockMovementService, request.ProductID, request.Quantity, constants.STOCK_MOVEMENT_TYPE_OUT, sell.ID, "Nueva venta")
	})

	if err != nil {
		return models.SellHistory{}, err
	}

	return sell, nil
}

func (s *sellHistoryService) DeleteSell(id uint64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txSellRepo := s.sellRepo.WithTx(tx)
		txProductStockService := s.productStockService.WithTx(tx)
		txStockMovementService := s.stockMovementService.WithTx(tx)

		sell, err := txSellRepo.FindByID(id)
		if err != nil {
			return err
		}

		if err := applyMovementFlow(txProductStockService, txStockMovementService, sell.ProductID, sell.Quantity, constants.STOCK_MOVEMENT_TYPE_IN, sell.ID, "Devolución de venta"); err != nil {
			return err
		}

		return txSellRepo.Delete(id)
	})
}
