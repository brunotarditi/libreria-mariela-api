package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	CreateMany(products []models.Product) (string, error)
	ExistsByCodeAndName(code, name string) (bool, error)
	ExistsBySku(sku string) (bool, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Model(&models.Product{}).Select("products.id, products.code, products.sku, products.name, products.profit_margin, products.description, Category.name AS category_name, Brand.name AS brand_name").Joins("Category").Joins("Brand").Find(&products).Error
	return products, err
}

func (r *productRepository) CreateMany(products []models.Product) (string, error) {
	err := r.db.CreateInBatches(&products, 100).Error
	return "Productos guardados con éxito", err
}

func (r *productRepository) ExistsByCodeAndName(code, name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("code = ? AND name = ?", code, name).Count(&count).Error
	return count > 0, err
}

func (r *productRepository) ExistsBySku(sku string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("sku = ?", sku).Count(&count).Error
	return count > 0, err
}
