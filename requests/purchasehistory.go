package requests

import (
	"fmt"
	"libreria/models"

	"gorm.io/gorm"
)

type PurchaseHistoryRequest struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	SupplierID uint    `json:"supplier_id" binding:"required"`
	Cost       float64 `json:"cost" binding:"required,gt=0"`
	Quantity   int     `json:"quantity" binding:"required,gt=0"`
}

func (r PurchaseHistoryRequest) ToModel() (models.PurchaseHistory, error) {
	return models.PurchaseHistory{
		ProductID:  r.ProductID,
		SupplierID: r.SupplierID,
		Cost:       r.Cost,
		Quantity:   r.Quantity,
	}, nil
}

func (r PurchaseHistoryRequest) UpdateModel(existing models.PurchaseHistory) (models.PurchaseHistory, error) {
	existing.ProductID = r.ProductID
	existing.SupplierID = r.SupplierID
	existing.Cost = r.Cost
	existing.Quantity = r.Quantity
	return existing, nil
}

func (r PurchaseHistoryRequest) Validate(db *gorm.DB) error {
	var product models.Product
	if err := db.First(&product, r.ProductID).Error; err != nil {
		return fmt.Errorf("producto con ID %d no encontrado", r.ProductID)
	}
	var supplier models.Supplier
	if err := db.First(&supplier, r.SupplierID).Error; err != nil {
		return fmt.Errorf("proveedor con ID %d no encontrado", r.SupplierID)
	}
	return nil
}
