package repositories

import (
	"fmt"
	"libreria/models"
	"strings"

	"gorm.io/gorm"
)

type SearchRepository interface {
	SearchProducts(query string, limit int) ([]models.Product, error)
	SearchBrands(query string, limit int) ([]models.Brand, error)
	SearchCategories(query string, limit int) ([]models.Category, error)
	SearchSuppliers(query string, limit int) ([]models.Supplier, error)
	SearchCustomers(query string, limit int) ([]models.Customer, error)
}

type searchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchProducts(query string, limit int) ([]models.Product, error) {
	var products []models.Product
	pattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Model(&models.Product{}).
		Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(description) LIKE ?", pattern, pattern, pattern, pattern).
		Limit(limit).
		Find(&products).Error
	return products, err
}

func (r *searchRepository) SearchBrands(query string, limit int) ([]models.Brand, error) {
	var brands []models.Brand
	pattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Model(&models.Brand{}).
		Where("LOWER(name) LIKE ?", pattern).
		Limit(limit).
		Find(&brands).Error
	return brands, err
}

func (r *searchRepository) SearchCategories(query string, limit int) ([]models.Category, error) {
	var categories []models.Category
	pattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Model(&models.Category{}).
		Where("LOWER(name) LIKE ?", pattern).
		Limit(limit).
		Find(&categories).Error
	return categories, err
}

func (r *searchRepository) SearchSuppliers(query string, limit int) ([]models.Supplier, error) {
	var suppliers []models.Supplier
	pattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Model(&models.Supplier{}).
		Where("LOWER(name) LIKE ? OR LOWER(contact_info) LIKE ?", pattern, pattern).
		Limit(limit).
		Find(&suppliers).Error
	return suppliers, err
}

func (r *searchRepository) SearchCustomers(query string, limit int) ([]models.Customer, error) {
	var customers []models.Customer
	pattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Model(&models.Customer{}).
		Where("LOWER(name) LIKE ? OR LOWER(contact_info) LIKE ?", pattern, pattern).
		Limit(limit).
		Find(&customers).Error
	return customers, err
}

// format helper for pattern
func toSearchPattern(query string) string {
	return fmt.Sprintf("%%%s%%", strings.ToLower(query))
}
