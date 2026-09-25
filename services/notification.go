package services

import (
	"fmt"
	"libreria/models"
	"libreria/repositories"
	"time"

	"gorm.io/gorm"
)

type NotificationService interface {
	GetNotifications(userID *uint, unreadOnly bool, limit int) ([]models.Notification, int64, error)
	MarkAsRead(id uint, userID *uint) error
	MarkAllAsRead(userID *uint) (int64, error)
	Delete(id uint, userID *uint) error
	ClearAll(userID *uint) (int64, error)
	CreateStockWarning(productID uint, productName string, sku string, currentStock int) error
	CreateWelcomeNotification(userID uint) error
	WithTx(tx *gorm.DB) NotificationService
}

type notificationService struct {
	db               *gorm.DB
	notificationRepo repositories.NotificationRepository
}

func NewNotificationService(db *gorm.DB, repo repositories.NotificationRepository) NotificationService {
	return &notificationService{
		db:               db,
		notificationRepo: repo,
	}
}

func (s *notificationService) WithTx(tx *gorm.DB) NotificationService {
	return &notificationService{
		db:               tx,
		notificationRepo: repositories.NewNotificationRepository(tx),
	}
}

func (s *notificationService) GetNotifications(userID *uint, unreadOnly bool, limit int) ([]models.Notification, int64, error) {
	notifications, err := s.notificationRepo.FindAll(userID, unreadOnly, limit)
	if err != nil {
		return nil, 0, err
	}

	unreadCount, err := s.notificationRepo.CountUnread(userID)
	if err != nil {
		return nil, 0, err
	}

	return notifications, unreadCount, nil
}

func (s *notificationService) MarkAsRead(id uint, userID *uint) error {
	return s.notificationRepo.MarkAsRead(id, userID)
}

func (s *notificationService) MarkAllAsRead(userID *uint) (int64, error) {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *notificationService) Delete(id uint, userID *uint) error {
	return s.notificationRepo.Delete(id, userID)
}

func (s *notificationService) ClearAll(userID *uint) (int64, error) {
	return s.notificationRepo.DeleteAll(userID)
}

func (s *notificationService) CreateStockWarning(productID uint, productName string, sku string, currentStock int) error {
	route := "/products"

	var title string
	var message string
	var nType string
	var icon string

	if currentStock <= 0 {
		title = fmt.Sprintf("Stock Agotado: %s", productName)
		message = fmt.Sprintf("El producto %s (SKU: %s) se ha quedado sin existencias (stock: 0).", productName, sku)
		nType = "alert"
		icon = "warning"
	} else {
		title = fmt.Sprintf("Stock Mínimo: %s", productName)
		message = fmt.Sprintf("El producto %s (SKU: %s) alcanzó el nivel de stock crítico (unidades restantes: %d).", productName, sku, currentStock)
		nType = "warning"
		icon = "inventory_2"
	}

	notification := models.Notification{
		UserID:    nil, // Notificación global para todos los administradores/vendedores
		Title:     title,
		Message:   message,
		Type:      nType,
		Icon:      icon,
		Route:     &route,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.notificationRepo.Create(&notification)
}

func (s *notificationService) CreateWelcomeNotification(userID uint) error {
	exists, err := s.notificationRepo.ExistsWelcome(userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	route := "/dashboard"
	uid := userID
	notification := models.Notification{
		UserID:    &uid,
		Title:     "¡Bienvenido a Librería Mariela!",
		Message:   "Tu cuenta ha sido conectada con éxito. Ya puedes consultar el catálogo, registrar ventas y administrar el inventario.",
		Type:      "info",
		Icon:      "store",
		Route:     &route,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.notificationRepo.Create(&notification)
}
