package requests

import (
	"fmt"
	"libreria/models"

	"gorm.io/gorm"
)

type ProductRequest struct {
	Name         string  `json:"name" binding:"required,min=1,max=65"`
	Code         string  `json:"code" binding:"required,min=1,max=20"`
	Sku          string  `json:"sku" binding:"max=20"`
	ProfitMargin float64 `json:"profit_margin" binding:"required,gt=0"`
	Description  string  `json:"description"`
	CategoryID   uint    `json:"category_id" binding:"required"`
	BrandID      uint    `json:"brand_id" binding:"required"`
}

func (r ProductRequest) ToModel() (models.Product, error) {
	return models.Product{
		Name:         r.Name,
		Code:         r.Code,
		Sku:          r.Sku,
		ProfitMargin: r.ProfitMargin,
		CategoryID:   r.CategoryID,
		Description:  r.Description,
		BrandID:      r.BrandID,
	}, nil
}

func (r ProductRequest) UpdateModel(existing models.Product) (models.Product, error) {
	existing.Name = r.Name
	existing.Code = r.Code
	existing.Sku = r.Sku
	existing.ProfitMargin = r.ProfitMargin
	existing.Description = r.Description
	existing.CategoryID = r.CategoryID
	existing.BrandID = r.BrandID
	return existing, nil
}

func (r ProductRequest) Validate(db *gorm.DB) error {
	var category models.Category
	if err := db.First(&category, r.CategoryID).Error; err != nil {
		return fmt.Errorf("categoría con ID %d no encontrada", r.CategoryID)
	}
	// Opcional: Validar que el nombre solo contenga letras
	// for _, c := range r.Name {
	// 	if c < 'A' || c > 'z' {
	// 		return fmt.Errorf("el nombre solo puede contener letras")
	// 	}
	// }
	return nil
}
