package db

import (
	"fmt"
	"libreria/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func ConnectDB() (db *gorm.DB) {

	// Validar variables de entorno
	requiredEnvVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Error: la variable de entorno %s no está definida", envVar)
		}
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})

	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	mysqlDB = db
	return mysqlDB
}

func AutoMigrate() {
	if mysqlDB == nil {
		log.Fatal("La base de datos no está inicializada")
	}
	mysqlDB.AutoMigrate(
		&models.Category{},
		&models.Brand{},
		&models.Customer{},
		&models.Supplier{},
		&models.Product{},
		&models.PurchaseHistory{},
		&models.SellHistory{},
		&models.PriceList{},
		&models.ProductStock{},
		&models.StockMovement{},
		&models.AuditLog{},
		&models.Role{},
		&models.User{},
		&models.UserRole{},
	)
}

func DisconnectDB() {
	if mysqlDB == nil {
		log.Fatal("La base de datos no está inicializada")
	}
	connect, err := mysqlDB.DB()
	if err != nil {
		log.Fatal("Error al obtener la conexión de la base de datos:", err)
	}
	if err := connect.Close(); err != nil {
		log.Fatal("Error al cerrar la conexión:", err)
	}
}

func GetDatabaseMySQL() *gorm.DB {
	if mysqlDB == nil {
		log.Fatal("La base de datos no está inicializada")
	}
	return mysqlDB
}
