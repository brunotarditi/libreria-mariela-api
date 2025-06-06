package main

import (
	"libreria/app"
	"libreria/common"
	"libreria/controllers"
	"libreria/middlewares"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"
	"libreria/security"
	"libreria/services"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.App) {

	// Security
	tokenManager, _ := security.NewPasetoManager()
	// Repositorios
	userRepo := repositories.NewUserRepositoryRepository(app.DB)
	purchaseRepo := repositories.NewPurchaseHistoryRepository(app.DB)
	sellRepo := repositories.NewSellHistoryRepository(app.DB)
	productStockRepo := repositories.NewProductStockRepository(app.DB)
	stockMovementRepo := repositories.NewStockMovementRepository(app.DB)
	roleRepo := repositories.NewRoleRepositoryRepository(app.DB)
	// Servicios
	userService := services.NewUserService(app.DB, userRepo, roleRepo, tokenManager)
	productStockService := services.NewProductStockService(app.DB, productStockRepo)
	stockMovementService := services.NewStockMovementService(app.DB, stockMovementRepo)
	purchaseService := services.NewPurchaseHistoryService(app.DB, purchaseRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)
	sellService := services.NewSellHistoryService(app.DB, sellRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)

	// Seed inicial
	if err := services.SeedInitialData(userRepo, roleRepo); err != nil {
		log.Fatalf("Error al hacer seed inicial: %v", err)
	}
	// Controladores
	userController := controllers.NewUserControllerController(userService)
	purchaseController := controllers.NewPurchaseHistoryController(purchaseService)
	sellController := controllers.NewSellHistoryControllerController(sellService)

	api := r.Group("/api/v1")

	public := api.Group("/users")
	{
		public.POST("/login", userController.Login())
		public.POST("/register", userController.UserRegister())

	}

	private := api.Group("/")
	private.Use(middlewares.AuthMiddleware(tokenManager, userService))
	private.Use(middlewares.AuditMiddleware(app.DB))

	{
		roles := private.Group("/roles")
		{
			roles.POST("/:id/assign", middlewares.RequireAnyRole("ROOT"), userController.AssignRole())
		}
		categories := private.Group("/categories")
		{
			ops := common.NewGormOperations[models.Category](app.DB)
			categories.GET("/", common.Get(ops))
			categories.GET("/:id", common.GetByID(ops))
			categories.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Category, requests.CategoryRequest](ops))
			categories.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Category, requests.CategoryRequest](ops))
			categories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}
		brands := private.Group("/brands")
		{
			ops := common.NewGormOperations[models.Brand](app.DB)
			brands.GET("/", common.Get(ops))
			brands.GET("/:id", common.GetByID(ops))
			brands.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Brand, requests.BrandRequest](ops))
			brands.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Brand, requests.BrandRequest](ops))
			brands.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}
		customers := private.Group("/customers")
		{
			ops := common.NewGormOperations[models.Customer](app.DB)
			customers.GET("/", common.Get(ops))
			customers.GET("/:id", common.GetByID(ops))
			customers.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Customer, requests.CustomerRequest](ops))
			customers.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Customer, requests.CustomerRequest](ops))
			customers.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}
		suppliers := private.Group("/suppliers")
		{
			ops := common.NewGormOperations[models.Supplier](app.DB)
			suppliers.GET("/", common.Get(ops))
			suppliers.GET("/:id", common.GetByID(ops))
			suppliers.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Supplier, requests.SupplierRequest](ops))
			suppliers.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Supplier, requests.SupplierRequest](ops))
			suppliers.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}
		products := private.Group("/products")
		{
			ops := common.NewGormOperations[models.Product](app.DB)
			products.GET("/", common.Get(ops))
			products.GET("/:id", common.GetByID(ops))
			products.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Product, requests.ProductRequest](ops))
			products.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Product, requests.ProductRequest](ops))
			products.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}
		purchaseHistories := private.Group("/purchases")
		{
			ops := common.NewGormOperations[models.PurchaseHistory](app.DB)
			purchaseHistories.GET("/", common.Get(ops))
			purchaseHistories.GET("/:id", common.GetByID(ops))
			purchaseHistories.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), purchaseController.CreatePurchaseHistory())
			purchaseHistories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), purchaseController.DeletePurchaseHistory())
		}
		sellHistories := private.Group("/sells")
		{
			ops := common.NewGormOperations[models.SellHistory](app.DB)
			sellHistories.GET("/", common.Get(ops))
			sellHistories.GET("/:id", common.GetByID(ops))
			sellHistories.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), sellController.CreateSellHistory())
			sellHistories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), sellController.DeleteSellHistory())
		}
		stocks := private.Group("/stocks")
		{
			ops := common.NewGormOperations[models.StockMovement](app.DB)
			stocks.GET("/", common.Get(ops))
			stocks.GET("/:id", common.GetByID(ops))
		}
		priceLists := private.Group("/prices")
		{
			ops := common.NewGormOperations[models.PriceList](app.DB)
			priceLists.GET("/", common.Get(ops))
			priceLists.GET("/:id", common.GetByID(ops))
			priceLists.POST("/", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.PriceList, requests.PriceListRequest](ops))
		}

	}
}
