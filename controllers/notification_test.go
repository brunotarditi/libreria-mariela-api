package controllers

import (
	"encoding/json"
	"libreria/models"
	"libreria/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mockNotificationService struct {
	notifications []models.Notification
}

func (m *mockNotificationService) GetNotifications(userID *uint, unreadOnly bool, limit int) ([]models.Notification, int64, error) {
	return m.notifications, 1, nil
}

func (m *mockNotificationService) MarkAsRead(id uint, userID *uint) error {
	return nil
}

func (m *mockNotificationService) MarkAllAsRead(userID *uint) (int64, error) {
	return 2, nil
}

func (m *mockNotificationService) Delete(id uint, userID *uint) error {
	return nil
}

func (m *mockNotificationService) ClearAll(userID *uint) (int64, error) {
	return 2, nil
}

func (m *mockNotificationService) CreateStockWarning(productID uint, productName string, sku string, currentStock int) error {
	return nil
}

func (m *mockNotificationService) CreateWelcomeNotification(userID uint) error {
	return nil
}

func (m *mockNotificationService) WithTx(tx *gorm.DB) services.NotificationService {
	return m
}

func TestNotificationController_Endpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockNotificationService{
		notifications: []models.Notification{
			{ID: 1, Title: "Aviso Stock", Type: "warning", IsRead: false},
			{ID: 2, Title: "Bienvenido", Type: "info", IsRead: true},
		},
	}
	ctrl := NewNotificationController(mockService)

	r := gin.New()
	r.GET("/notifications", ctrl.GetNotifications)
	r.PATCH("/notifications/:id/read", ctrl.MarkAsRead)
	r.PATCH("/notifications/read-all", ctrl.MarkAllAsRead)
	r.DELETE("/notifications/:id", ctrl.Delete)
	r.DELETE("/notifications", ctrl.ClearAll)

	// 1. GET /notifications
	t.Run("GET /notifications", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/notifications?unread=true&limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		if res["unread_count"] != float64(1) {
			t.Errorf("expected unread_count 1, got %v", res["unread_count"])
		}
	})

	// 2. PATCH /notifications/1/read
	t.Run("PATCH /notifications/:id/read", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/notifications/1/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// 3. PATCH /notifications/read-all
	t.Run("PATCH /notifications/read-all", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// 4. DELETE /notifications/1
	t.Run("DELETE /notifications/:id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/notifications/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// 5. DELETE /notifications
	t.Run("DELETE /notifications", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/notifications", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})
}
