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
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.App) {

	// Security
	tokenManager, _ := security.NewPasetoManager()
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
	roleRepo := repositories.NewRoleRepositoryRepository(app.DB)
	sellRepo := repositories.NewSellHistoryRepository(app.DB)
	stockMovementRepo := repositories.NewStockMovementRepository(app.DB)
	userRepo := repositories.NewUserRepositoryRepository(app.DB)
	emailVerificationRepo := repositories.NewEmailVerificationRepositoryRepository(app.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(app.DB)
	// Servicios
	productService := services.NewProductService(app.DB, productRepo, categoryOps, brandOps)
	productStockService := services.NewProductStockService(app.DB, productStockRepo)
	stockMovementService := services.NewStockMovementService(app.DB, stockMovementRepo)
	purchaseService := services.NewPurchaseHistoryService(app.DB, purchaseRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)
	sellService := services.NewSellHistoryService(app.DB, sellRepo, productStockRepo, stockMovementRepo, productStockService, stockMovementService)
	userService := services.NewUserService(app.DB, userRepo, roleRepo, tokenManager, emailVerificationRepo, passwordResetRepo)
	dashboardService := services.NewDashboardService(app.DB, dashboardRepo, supplierOps, customerdOps, productOps)
	budgetService := services.NewBudgetService(app.DB)

	// Seed inicial
	if err := services.SeedInitialData(userRepo, roleRepo); err != nil {
		log.Fatalf("Error al hacer seed inicial: %v", err)
	}
	// Controladores
	productController := controllers.NewProductController(productService)
	purchaseController := controllers.NewPurchaseHistoryController(purchaseService)
	sellController := controllers.NewSellHistoryControllerController(sellService)
	userController := controllers.NewUserControllerController(userService)
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

	verify := router.Group("/verify-email")
	{
		verify.GET("", userController.VerifyEmail())
	}

	resetpass := router.Group("/reset-password")
	{
		resetpass.GET("", userController.ForgotPassword())
		resetpass.POST("", userController.ResetPassword())
	}

	public := router.Group("/auth")
	{
		public.POST("/login", userController.Login())
		public.POST("/register", userController.Register())
		public.GET("/verify-token", userController.Register())
	}

	private := router.Group("/")
	private.Use(middlewares.AuthMiddleware(tokenManager, userService))
	private.Use(middlewares.AuditMiddleware(app.DB))

	{
		users := private.Group("/users")
		{
			ops := common.NewGormOperations[models.User](app.DB)
			users.GET("", userController.FindAll())
			users.PATCH("/:id", common.Update[models.User, requests.UserRequest](ops))
		}

		roles := private.Group("/roles")
		{
			roles.POST("/:id/assign", middlewares.RequireAnyRole("ROOT"), userController.AssignRole())
		}

		categories := private.Group("/categories")
		{
			categories.GET("", common.Get(categoryOps))
			categories.GET("/:id", common.GetByID(categoryOps))
			categories.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Category, requests.CategoryRequest](categoryOps))
			categories.POST("/list", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.CreateMany[models.Category, requests.CategoryRequestArray](categoryOps))
			categories.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Category, requests.CategoryRequest](categoryOps))
			categories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(categoryOps))
		}
		brands := private.Group("/brands")
		{
			brands.GET("", common.Get(brandOps))
			brands.GET("/:id", common.GetByID(brandOps))
			brands.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Brand, requests.BrandRequest](brandOps))
			brands.POST("/list", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.CreateMany[models.Brand, requests.BrandRequestArray](brandOps))
			brands.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Brand, requests.BrandRequest](brandOps))
			brands.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(brandOps))
		}
		customers := private.Group("/customers")
		{
			customers.GET("", common.Get(customerdOps))
			customers.GET("/:id", common.GetByID(customerdOps))
			customers.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Customer, requests.CustomerRequest](customerdOps))
			customers.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Customer, requests.CustomerRequest](customerdOps))
			customers.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(customerdOps))
		}
		suppliers := private.Group("/suppliers")
		{
			suppliers.GET("", common.Get(supplierOps))
			suppliers.GET("/:id", common.GetByID(supplierOps))
			suppliers.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Supplier, requests.SupplierRequest](supplierOps))
			suppliers.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Supplier, requests.SupplierRequest](supplierOps))
			suppliers.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(supplierOps))
		}
		products := private.Group("/products")
		{
			products.GET("/:id", common.GetByID(productOps))
			products.GET("", productController.FindAllWithCategoriesAndBrands())
			products.GET("/export", middlewares.RequireAnyRole("WRITE", "READ", "ADMIN", "ROOT"), productController.GetExport())
			products.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Product, requests.ProductRequest](productOps))
			products.POST("/import", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), productController.ImportFromExcel())
			products.PUT("/:id", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Update[models.Product, requests.ProductRequest](productOps))
			products.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(productOps))
		}
		purchaseHistories := private.Group("/purchases")
		{
			ops := common.NewGormOperations[models.PurchaseHistory](app.DB)
			purchaseHistories.GET("", common.Get(ops))
			purchaseHistories.GET("/:id", common.GetByID(ops))
			purchaseHistories.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), purchaseController.CreatePurchaseHistory())
			purchaseHistories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), purchaseController.DeletePurchaseHistory())
		}
		sellHistories := private.Group("/sells")
		{
			ops := common.NewGormOperations[models.SellHistory](app.DB)
			sellHistories.GET("", common.Get(ops))
			sellHistories.GET("/:id", common.GetByID(ops))
			sellHistories.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), sellController.CreateSellHistory())
			sellHistories.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), sellController.DeleteSellHistory())
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
			priceLists.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.PriceList, requests.PriceListRequest](ops))
		}
		dashboard := private.Group("/dashboard")
		{
			dashboard.GET("", dashboardController.GetData())
		}

		budges := private.Group("/budges")
		{
			ops := common.NewGormOperations[models.Budget](app.DB)
			budges.POST("", middlewares.RequireAnyRole("WRITE", "ADMIN", "ROOT"), common.Create[models.Budget, requests.BudgetRequest](ops))
			budges.DELETE("/:id", middlewares.RequireAnyRole("ADMIN", "ROOT"), common.Delete(ops))
		}

	}
}
