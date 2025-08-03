package repositories

import (
	"libreria/models"
	"libreria/responses"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]responses.ProductResponse, error)
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

func (r *productRepository) FindAll() ([]responses.ProductResponse, error) {
	var products []responses.ProductResponse
	err := r.db.Model(&models.Product{}).
		Select("products.id, products.code, products.sku, products.name, products.profit_margin, products.description, category.name AS category_name, brand.name AS brand_name").
		Joins("LEFT JOIN categories category ON category.id = products.category_id").
		Joins("LEFT JOIN brands brand ON brand.id = products.brand_id").
		Find(&products).Error
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
