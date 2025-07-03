package controllers

import (
	"libreria/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) FindAllWithCategoriesAndBrands() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		products, err := c.service.GetAllProductsWithCategoriesAndBrands()
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, products)
	}
}

func (c *ProductController) GetExport() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		file, err := c.service.ExportToExcel()
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		filePath := "assets/templates/products.xlsx"
		if err := file.SaveAs(filePath); err != nil {
			ctx.JSON(500, gin.H{"error": "error al guardar el archivo Excel"})
			return
		}

		ctx.FileAttachment(filePath, "products.xlsx")
	}
}

func (c *ProductController) ImportFromExcel() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Archivo requerido"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(500, gin.H{"error": "No se pudo abrir el archivo"})
			return
		}
		defer file.Close()

		err = c.service.ImportFromExcel(file)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "Importación exitosa"})
	}
}
