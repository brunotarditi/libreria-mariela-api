package controllers

import (
	"libreria/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service services.DashboardService
}

func NewDashboardController(service services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (c *DashboardController) GetData() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		response, err := c.service.GetData()
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}
