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

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txPurchaseRepo := s.purchaseRepo.WithTx(tx)
		txProductStockService := s.productStockService.WithTx(tx)
		txStockMovementService := s.stockMovementService.WithTx(tx)

		if err := txPurchaseRepo.Create(&purchase); err != nil {
			return err
		}

		return applyMovementFlow(txProductStockService, txStockMovementService, request.ProductID, request.Quantity, constants.STOCK_MOVEMENT_TYPE_IN, purchase.ID, "Nueva compra")
	})

	if err != nil {
		return models.PurchaseHistory{}, err
	}

	return purchase, nil
}

func (s *purchaseHistoryService) DeletePurchase(id uint64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txPurchaseRepo := s.purchaseRepo.WithTx(tx)
		txProductStockService := s.productStockService.WithTx(tx)
		txStockMovementService := s.stockMovementService.WithTx(tx)

		purchase, err := txPurchaseRepo.FindByID(id)
		if err != nil {
			return err
		}

		if err := applyMovementFlow(txProductStockService, txStockMovementService, purchase.ProductID, purchase.Quantity, constants.STOCK_MOVEMENT_TYPE_OUT, purchase.ID, "Devolución de compra"); err != nil {
			return err
		}

		return txPurchaseRepo.Delete(id)
	})
}
