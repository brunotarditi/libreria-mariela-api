package controllers

import (
	"libreria/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	service services.SearchService
}

func NewSearchController(service services.SearchService) *SearchController {
	return &SearchController{service: service}
}

func (ctrl *SearchController) Search(c *gin.Context) {
	q := c.Query("q")
	result, err := ctrl.service.Search(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
