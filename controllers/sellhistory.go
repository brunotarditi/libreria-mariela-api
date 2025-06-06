package controllers

import (
	"libreria/requests"
	"libreria/services"

	"github.com/gin-gonic/gin"
)

type SellHistoryController struct {
	service services.SellHistoryService
}

func NewSellHistoryControllerController(service services.SellHistoryService) *SellHistoryController {
	return &SellHistoryController{service: service}
}

func (c *SellHistoryController) CreateSellHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.SellHistoryRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		sell, err := c.service.CreateSell(request)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, sell)
	}
}

func (c *SellHistoryController) DeleteSellHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := c.service.DeleteSell(id); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(204, nil)
	}
}
