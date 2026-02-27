package main

import (
	"libreria/app"
	"libreria/common"
	"libreria/controllers"
	"libreria/middlewares"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"
	"libreria/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.App) {

	// Common operations
	categoryOps := common.NewGormOperations[models.Category](app.DB)
	brandOps := common.NewGormOperations[models.Brand](app.DB)
	productOps := common.NewGormOperations[models.Product](app.DB)
	supplierOps := common.NewGormOperations[models.Supplier](app.DB)
	customerdOps := common.NewGormOperations[models.Customer](app.DB)
	// Repositorios
	dashboardRepo := repositories.NewDashboardRepository(app.DB)
	productRepo := repositories.NewProductRepository(app.DB)
	productStockRepo := repositories.NewProductStockRepository(app.DB)
	purchaseRepo := repositories.NewPurchaseHistoryRepository(app.DB)
	sellRepo := repositories.NewSellHistoryRepository(app.DB)
	stockMovementRepo := repositories.NewStockMovementRepository(app.DB)
	// Servicios
	productService := services.NewProductService(app.DB, productRepo, categoryOps, brandOps)
	productStockService := services.NewProductStockService(app.DB, productStockRepo)
	stockMovementService := services.NewStockMovementService(app.DB, stockMovementRepo)
	purchaseService := services.NewPurchaseHistoryService(app.DB, purchaseRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)
	sellService := services.NewSellHistoryService(app.DB, sellRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)
	dashboardService := services.NewDashboardService(app.DB, dashboardRepo, supplierOps, customerdOps, productOps)
	budgetService := services.NewBudgetService(app.DB)

	// Controladores
	productController := controllers.NewProductController(productService)
	purchaseController := controllers.NewPurchaseHistoryController(purchaseService)
	sellController := controllers.NewSellHistoryControllerController(sellService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	budgetController := controllers.NewBudgetController(budgetService)

	router := r.Group("/api/v1")

	budget := router.Group("/budget")
	{
		budget.GET("", budgetController.GetBudget())
	}

	health := router.Group("/healthy")
	{
		health.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, map[string]interface{}{"message": "success", "api": "libreria mariela", "status": "http.StatusOK", "description": "This is the REST API for the Libreria Mariela store"})
		})
	}

	private := router.Group("/")
	private.Use(middlewares.AuditMiddleware(app.DB))

	{

		categories := private.Group("/categories")
		{
			categories.GET("", common.Get(categoryOps))
			categories.GET("/:id", common.GetByID(categoryOps))
			categories.POST("", common.Create[models.Category, requests.CategoryRequest](categoryOps))
			categories.POST("/list", common.CreateMany[models.Category, requests.CategoryRequestArray](categoryOps))
			categories.PUT("/:id", common.Update[models.Category, requests.CategoryRequest](categoryOps))
			categories.DELETE("/:id", common.Delete(categoryOps))
		}
		brands := private.Group("/brands")
		{
			brands.GET("", common.Get(brandOps))
			brands.GET("/:id", common.GetByID(brandOps))
			brands.POST("", common.Create[models.Brand, requests.BrandRequest](brandOps))
			brands.POST("/list", common.CreateMany[models.Brand, requests.BrandRequestArray](brandOps))
			brands.PUT("/:id", common.Update[models.Brand, requests.BrandRequest](brandOps))
			brands.DELETE("/:id", common.Delete(brandOps))
		}
		customers := private.Group("/customers")
		{
			customers.GET("", common.Get(customerdOps))
			customers.GET("/:id", common.GetByID(customerdOps))
			customers.POST("", common.Create[models.Customer, requests.CustomerRequest](customerdOps))
			customers.PUT("/:id", common.Update[models.Customer, requests.CustomerRequest](customerdOps))
			customers.DELETE("/:id", common.Delete(customerdOps))
		}
		suppliers := private.Group("/suppliers")
		{
			suppliers.GET("", common.Get(supplierOps))
			suppliers.GET("/:id", common.GetByID(supplierOps))
			suppliers.POST("", common.Create[models.Supplier, requests.SupplierRequest](supplierOps))
			suppliers.PUT("/:id", common.Update[models.Supplier, requests.SupplierRequest](supplierOps))
			suppliers.DELETE("/:id", common.Delete(supplierOps))
		}
		products := private.Group("/products")
		{
			products.GET("/:id", common.GetByID(productOps))
			products.GET("", productController.FindAllWithCategoriesAndBrands())
			products.GET("/export", productController.GetExport())
			products.POST("", common.Create[models.Product, requests.ProductRequest](productOps))
			products.POST("/import", productController.ImportFromExcel())
			products.PUT("/:id", common.Update[models.Product, requests.ProductRequest](productOps))
			products.DELETE("/:id", common.Delete(productOps))
		}
		purchaseHistories := private.Group("/purchases")
		{
			ops := common.NewGormOperations[models.PurchaseHistory](app.DB)
			purchaseHistories.GET("", common.Get(ops))
			purchaseHistories.GET("/:id", common.GetByID(ops))
			purchaseHistories.POST("", purchaseController.CreatePurchaseHistory())
			purchaseHistories.DELETE("/:id", purchaseController.DeletePurchaseHistory())
		}
		sellHistories := private.Group("/sells")
		{
			ops := common.NewGormOperations[models.SellHistory](app.DB)
			sellHistories.GET("", common.Get(ops))
			sellHistories.GET("/:id", common.GetByID(ops))
			sellHistories.POST("", sellController.CreateSellHistory())
			sellHistories.DELETE("/:id", sellController.DeleteSellHistory())
		}
		stocks := private.Group("/stocks")
		{
			ops := common.NewGormOperations[models.StockMovement](app.DB)
			stocks.GET("", common.Get(ops))
			stocks.GET("/:id", common.GetByID(ops))
		}
		priceLists := private.Group("/prices")
		{
			ops := common.NewGormOperations[models.PriceList](app.DB)
			priceLists.GET("", common.Get(ops))
			priceLists.GET("/:id", common.GetByID(ops))
			priceLists.POST("", common.Create[models.PriceList, requests.PriceListRequest](ops))
		}
		dashboard := private.Group("/dashboard")
		{
			dashboard.GET("", dashboardController.GetData())
		}

		budges := private.Group("/budges")
		{
			ops := common.NewGormOperations[models.Budget](app.DB)
			budges.POST("", common.Create[models.Budget, requests.BudgetRequest](ops))
			budges.DELETE("/:id", common.Delete(ops))
		}

	}
}
