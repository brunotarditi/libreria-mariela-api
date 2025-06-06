package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type ProductStockRepository interface {
	Create(stock *models.ProductStock) error
	FindByID(id uint) (models.ProductStock, error)
	Update(productstock *models.ProductStock) error
}

type productStockRepository struct {
	db *gorm.DB
}

func NewProductStockRepository(db *gorm.DB) ProductStockRepository {
	return &productStockRepository{db: db}
}

func (r *productStockRepository) Create(productstock *models.ProductStock) error {
	return r.db.Create(productstock).Error
}

func (r *productStockRepository) FindByID(id uint) (models.ProductStock, error) {
	var stock models.ProductStock
	err := r.db.First(&stock, id).Error
	return stock, err
}

func (r *productStockRepository) Update(productstock *models.ProductStock) error {
	return r.db.Save(productstock).Error
}
