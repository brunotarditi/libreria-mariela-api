package requests

import (
	"libreria/models"
	"time"
)

type PriceListRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

func (r PriceListRequest) ToModel() (models.PriceList, error) {
	priceList := models.PriceList{
		ProductID:   r.ProductID,
		Price:       r.Price,
		EffectiveAt: time.Now(),
		IsActive:    true,
	}

	return priceList, nil
}

func (r PriceListRequest) UpdateModel(existing models.PriceList) (models.PriceList, error) {
	existing.Price = r.Price
	existing.EffectiveAt = time.Now()
	existing.IsActive = true
	return existing, nil
}
