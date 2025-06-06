package controllers

import (
	"fmt"
	"libreria/requests"
	"libreria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseHistoryController struct {
	service services.PurchaseHistoryService
}

func NewPurchaseHistoryController(service services.PurchaseHistoryService) *PurchaseHistoryController {
	return &PurchaseHistoryController{service: service}
}

func (c *PurchaseHistoryController) CreatePurchaseHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.PurchaseHistoryRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		purchase, err := c.service.CreatePurchase(request)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, purchase)
	}
}

func (c *PurchaseHistoryController) DeletePurchaseHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		purchaseHistoryID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}
		if err := c.service.DeletePurchase(purchaseHistoryID); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(204, nil)
	}
}
