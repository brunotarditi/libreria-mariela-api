package controllers

import (
	"libreria/requests"
	"libreria/services"
	"net/http"
	"strconv"

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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sell, err := c.service.CreateSell(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, sell)
	}
}

func (c *SellHistoryController) DeleteSellHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		sellHistoryID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido: debe ser numérico"})
			return
		}
		if err := c.service.DeleteSell(sellHistoryID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(204, nil)
	}
}
