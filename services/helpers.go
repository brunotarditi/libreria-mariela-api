package services

import (
	"libreria/constants"
	"libreria/models"
)

func applyMovementFlow(productStockService ProductStockService, stockMovementService StockMovementService, productID uint, qty int, movType constants.StockMovementType, refID uint, note string) error {

	productStock := models.ProductStock{ProductID: productID, Quantity: qty}
	if err := productStockService.ApplyMovement(productStock, int(movType)); err != nil {
		return err
	}

	stockMovement := models.StockMovement{
		ProductID:    productID,
		Quantity:     qty,
		MovementType: int(movType),
		ReferenceID:  &refID,
		Note:         note,
	}
	if err := stockMovementService.ApplyMovement(stockMovement); err != nil {
		return err
	}

	return nil
}
