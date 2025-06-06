package main

import (
	"libreria/app"
	"libreria/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env:", err)
	}

	dbInstance := db.ConnectDB()
	defer db.DisconnectDB()
	db.AutoMigrate()
	appInstance := app.NewApp(dbInstance)

	r := gin.New()

	// Agregar middlewares esenciales manualmente
	r.Use(gin.Recovery())

	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}

	r.SetTrustedProxies([]string{"127.0.0.1"})
	SetupRoutes(r, appInstance)

	r.Run(":8080")
}
