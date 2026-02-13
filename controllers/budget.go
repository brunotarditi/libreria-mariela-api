package controllers

import (
	"libreria/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	service services.BudgetService
}

func NewBudgetController(service services.BudgetService) *BudgetController {
	return &BudgetController{service: service}
}

func (c *BudgetController) GetBudget() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Generar el PDF
		pdfBytes, err := c.service.GeneratePDF()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
			return
		}

		ctx.Header("Content-Type", "application/pdf")
		ctx.Header("Content-Disposition", "attachment; filename=presupuesto.pdf")
		ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}
