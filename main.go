package main

import (
	"libreria/app"
	"libreria/db"
	"libreria/middlewares"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar variables de entorno desde el archivo .env
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No se pudo cargar .env, probablemente estés en producción")
		}
	}

	dbInstance := db.ConnectDB()
	defer db.DisconnectDB()
	db.AutoMigrate()
	appInstance := app.NewApp(dbInstance)

	r := gin.New()

	// Agregar middlewares esenciales manualmente
	r.Use(gin.Recovery())
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	r.Use(middlewares.CORSMiddleware(origins))
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}

	r.SetTrustedProxies([]string{"127.0.0.1"})
	SetupRoutes(r, appInstance)

	r.Run(":8080")
}
