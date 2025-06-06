package requests

import (
	"fmt"
	"libreria/models"

	"gorm.io/gorm"
)

type SellHistoryRequest struct {
	ProductID  uint `json:"product_id" binding:"required"`
	CustomerID uint `json:"customer_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,gt=0"`
}

func (r SellHistoryRequest) ToModel() (models.SellHistory, error) {
	return models.SellHistory{
		ProductID:   r.ProductID,
		CustomerID:  r.CustomerID,
		Price:       0,
		Quantity:    r.Quantity,
		AverageCost: 0,
	}, nil
}

func (r SellHistoryRequest) Validate(db *gorm.DB) error {
	var product models.Product
	if err := db.First(&product, r.ProductID).Error; err != nil {
		return fmt.Errorf("producto con ID %d no encontrado", r.ProductID)
	}
	var customer models.Customer
	if err := db.First(&customer, r.CustomerID).Error; err != nil {
		return fmt.Errorf("cliente con ID %d no encontrado", r.CustomerID)
	}

	return nil
}
