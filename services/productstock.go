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
	WithTx(tx *gorm.DB) ProductStockService
}

type productStockService struct {
	db                  *gorm.DB
	productStockRepo    repositories.ProductStockRepository
	notificationService NotificationService
}

func NewProductStockService(db *gorm.DB, stockRepo repositories.ProductStockRepository, notificationService ...NotificationService) ProductStockService {
	var notifService NotificationService
	if len(notificationService) > 0 {
		notifService = notificationService[0]
	}
	return &productStockService{
		db:                  db,
		productStockRepo:    stockRepo,
		notificationService: notifService,
	}
}

func (s *productStockService) WithTx(tx *gorm.DB) ProductStockService {
	var notifService NotificationService
	if s.notificationService != nil {
		notifService = s.notificationService.WithTx(tx)
	}
	return &productStockService{
		db:                  tx,
		productStockRepo:    s.productStockRepo.WithTx(tx),
		notificationService: notifService,
	}
}

func (s *productStockService) ApplyMovement(productStock models.ProductStock, movementType int) error {
	stockExist, err := s.productStockRepo.FindByID(productStock.ProductID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		if err := s.productStockRepo.Create(&productStock); err != nil {
			return err
		}
		return nil
	}

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
		return err
	}

	// Trigger automático: alerta cuando el stock llega a 0 o al stock mínimo (<= 5 unidades)
	if movementType == int(constants.STOCK_MOVEMENT_TYPE_OUT) && stockExist.Quantity <= 5 && s.notificationService != nil {
		var product models.Product
		if err := s.db.Select("id, name, sku").First(&product, stockExist.ProductID).Error; err == nil {
			_ = s.notificationService.CreateStockWarning(product.ID, product.Name, product.Sku, stockExist.Quantity)
		}
	}

	return nil
}
