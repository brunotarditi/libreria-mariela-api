package utils

import (
	"libreria/models"

	"gorm.io/gorm"
)

func CalculateAverageCostAndStock(db *gorm.DB, productID uint) (float64, int64, error) {
	var productStock models.ProductStock
	if err := db.Where("product_id = ?", productID).First(&productStock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	stock := int64(productStock.Quantity)
	if stock <= 0 {
		return 0, 0, nil
	}
	var purchaseHistories []models.PurchaseHistory
	if err := db.Where("product_id = ?", productID).Order("created_at asc").Find(&purchaseHistories).Error; err != nil {
		return 0, 0, err
	}

	var totalCost float64
	remaining := stock
	for _, ph := range purchaseHistories {
		if remaining <= 0 {
			break
		}
		quantity := int64(ph.Quantity)
		if quantity > remaining {
			quantity = remaining
		}
		totalCost += float64(quantity) * ph.Cost
		remaining -= quantity
	}

	if stock == 0 {
		return 0, 0, nil
	}
	averageCost := totalCost / float64(stock)

	return averageCost, stock, nil
}
