package services

import (
	"libreria/constants"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"

	"gorm.io/gorm"
)

type PurchaseHistoryService interface {
	CreatePurchase(request requests.PurchaseHistoryRequest) (models.PurchaseHistory, error)
	DeletePurchase(id uint64) error
}

type purchaseHistoryService struct {
	db                   *gorm.DB
	purchaseRepo         repositories.PurchaseHistoryRepository
	productStockRepo     repositories.ProductStockRepository
	stockMovementRepo    repositories.StockMovementRepository
	productStockService  ProductStockService
	stockMovementService StockMovementService
}

func NewPurchaseHistoryService(db *gorm.DB, purchaseRepo repositories.PurchaseHistoryRepository, productStockRepo repositories.ProductStockRepository, stockMovementRepo repositories.StockMovementRepository, productStockService ProductStockService, stockMovementService StockMovementService) PurchaseHistoryService {
	return &purchaseHistoryService{
		db:                   db,
		purchaseRepo:         purchaseRepo,
		productStockRepo:     productStockRepo,
		stockMovementRepo:    stockMovementRepo,
		productStockService:  productStockService,
		stockMovementService: stockMovementService,
	}
}

func (s *purchaseHistoryService) CreatePurchase(request requests.PurchaseHistoryRequest) (models.PurchaseHistory, error) {
	if err := request.Validate(s.db); err != nil {
		return models.PurchaseHistory{}, err
	}

	purchase, err := request.ToModel()
	if err != nil {
		return models.PurchaseHistory{}, err
	}

	tx := s.db.Begin()
	if err := s.purchaseRepo.Create(&purchase); err != nil {
		tx.Rollback()
		return models.PurchaseHistory{}, err
	}

	if err := applyMovementFlow(s.productStockService, s.stockMovementService, request.ProductID, request.Quantity, constants.STOCK_MOVEMENT_TYPE_IN, purchase.ID, "Nueva compra"); err != nil {
		tx.Rollback()
		return models.PurchaseHistory{}, err
	}

	tx.Commit()
	return purchase, nil
}

func (s *purchaseHistoryService) DeletePurchase(id uint64) error {
	tx := s.db.Begin()
	purchase, err := s.purchaseRepo.FindByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := applyMovementFlow(s.productStockService, s.stockMovementService, purchase.ProductID, purchase.Quantity, constants.STOCK_MOVEMENT_TYPE_OUT, purchase.ID, "Devolución de compra"); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.purchaseRepo.Delete(id); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
