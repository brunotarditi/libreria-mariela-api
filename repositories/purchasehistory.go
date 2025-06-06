package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type PurchaseHistoryRepository interface {
	Create(purchase *models.PurchaseHistory) error
	FindByID(purchaseHistoryID uint64) (models.PurchaseHistory, error)
	Delete(purchaseHistoryID uint64) error
}

type purchaseHistoryRepository struct {
	db *gorm.DB
}

func NewPurchaseHistoryRepository(db *gorm.DB) PurchaseHistoryRepository {
	return &purchaseHistoryRepository{db: db}
}

func (r *purchaseHistoryRepository) Create(purchase *models.PurchaseHistory) error {
	return r.db.Create(purchase).Error
}

func (r *purchaseHistoryRepository) FindByID(purchaseHistoryID uint64) (models.PurchaseHistory, error) {
	var purchase models.PurchaseHistory
	err := r.db.First(&purchase, purchaseHistoryID).Error
	return purchase, err
}

func (r *purchaseHistoryRepository) Delete(purchaseHistoryID uint64) error {
	return r.db.Delete(&models.PurchaseHistory{}, purchaseHistoryID).Error
}
