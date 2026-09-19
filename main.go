package main

import (
	"libreria/app"
	"libreria/db"
	"libreria/middlewares"
	"log"
	"os"
	"strings"

	"github.com/brunotarditi/peak-auth/sdk/go"
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

	peakAuthURL := os.Getenv("PEAK_AUTH_URL")
	if peakAuthURL == "" {
		peakAuthURL = "http://localhost:8080"
	}
	clientID := os.Getenv("PEAK_AUTH_CLIENT_ID")
	if clientID == "" {
		clientID = "libreria-mariela"
	}
	clientSecret := os.Getenv("PEAK_AUTH_CLIENT_SECRET")

	peakAuthClient, err := peakauth.New(peakauth.Config{
		IssuerURL:    peakAuthURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Fatalf("Error inicializando Peak Auth SDK: %v", err)
	}

	appInstance := app.NewApp(dbInstance, peakAuthClient)

	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 MB límite para subida de archivos
	// Agregar middlewares esenciales manualmente
	r.Use(gin.Recovery())
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	r.Use(middlewares.CORSMiddleware(origins))
	r.Use(gin.Logger())

	r.SetTrustedProxies(nil)
	SetupRoutes(r, appInstance)

	r.Run(":8080")
}
