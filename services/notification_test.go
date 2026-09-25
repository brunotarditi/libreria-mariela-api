package services

import (
	"libreria/models"
	"testing"
)

type mockNotificationRepo struct {
	notifications []models.Notification
}

func (m *mockNotificationRepo) FindAll(userID *uint, unreadOnly bool, limit int) ([]models.Notification, error) {
	var result []models.Notification
	for _, n := range m.notifications {
		if unreadOnly && n.IsRead {
			continue
		}
		if userID != nil && n.UserID != nil && *n.UserID != *userID {
			continue
		}
		result = append(result, n)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (m *mockNotificationRepo) CountUnread(userID *uint) (int64, error) {
	var count int64
	for _, n := range m.notifications {
		if !n.IsRead {
			if userID == nil || n.UserID == nil || *n.UserID == *userID {
				count++
			}
		}
	}
	return count, nil
}

func (m *mockNotificationRepo) FindByID(id uint) (models.Notification, error) {
	for _, n := range m.notifications {
		if n.ID == id {
			return n, nil
		}
	}
	return models.Notification{}, nil
}

func (m *mockNotificationRepo) MarkAsRead(id uint, userID *uint) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications[i].IsRead = true
			return nil
		}
	}
	return nil
}

func (m *mockNotificationRepo) MarkAllAsRead(userID *uint) (int64, error) {
	var count int64
	for i, n := range m.notifications {
		if !n.IsRead {
			m.notifications[i].IsRead = true
			count++
		}
	}
	return count, nil
}

func (m *mockNotificationRepo) Delete(id uint, userID *uint) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications = append(m.notifications[:i], m.notifications[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockNotificationRepo) DeleteAll(userID *uint) (int64, error) {
	count := int64(len(m.notifications))
	m.notifications = []models.Notification{}
	return count, nil
}

func (m *mockNotificationRepo) Create(notification *models.Notification) error {
	notification.ID = uint(len(m.notifications) + 1)
	m.notifications = append(m.notifications, *notification)
	return nil
}

func (m *mockNotificationRepo) ExistsUnreadByRoute(route string) (bool, error) {
	for _, n := range m.notifications {
		if n.Route != nil && *n.Route == route && !n.IsRead {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockNotificationRepo) ExistsWelcome(userID uint) (bool, error) {
	for _, n := range m.notifications {
		if n.UserID != nil && *n.UserID == userID && n.Icon == "store" {
			return true, nil
		}
	}
	return false, nil
}

func TestNotificationService_GetNotifications(t *testing.T) {
	repo := &mockNotificationRepo{
		notifications: []models.Notification{
			{ID: 1, Title: "Aviso 1", IsRead: false},
			{ID: 2, Title: "Aviso 2", IsRead: true},
		},
	}
	service := NewNotificationService(nil, repo)

	// All notifications
	notifs, unread, err := service.GetNotifications(nil, false, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifs) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(notifs))
	}
	if unread != 1 {
		t.Errorf("expected 1 unread, got %d", unread)
	}

	// Unread only
	notifs, _, err = service.GetNotifications(nil, true, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifs) != 1 {
		t.Errorf("expected 1 unread notification, got %d", len(notifs))
	}
}

func TestNotificationService_CreateStockWarning(t *testing.T) {
	repo := &mockNotificationRepo{}
	service := NewNotificationService(nil, repo)

	// Stock 0
	err := service.CreateStockWarning(1, "Lapicera", "SKU123", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.notifications) != 1 {
		t.Fatalf("expected 1 notification created, got %d", len(repo.notifications))
	}
	if repo.notifications[0].Type != "alert" || repo.notifications[0].Icon != "warning" {
		t.Errorf("expected alert/warning for stock 0, got type=%s icon=%s", repo.notifications[0].Type, repo.notifications[0].Icon)
	}

	// Stock critical (e.g. 3 units)
	err = service.CreateStockWarning(2, "Cuaderno", "SKU456", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.notifications) != 2 {
		t.Fatalf("expected 2 notifications created, got %d", len(repo.notifications))
	}
	if repo.notifications[1].Type != "warning" || repo.notifications[1].Icon != "inventory_2" {
		t.Errorf("expected warning/inventory_2 for stock 3, got type=%s icon=%s", repo.notifications[1].Type, repo.notifications[1].Icon)
	}
}

func TestNotificationService_CreateWelcomeNotification(t *testing.T) {
	repo := &mockNotificationRepo{}
	service := NewNotificationService(nil, repo)

	err := service.CreateWelcomeNotification(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.notifications) != 1 {
		t.Fatalf("expected 1 notification created, got %d", len(repo.notifications))
	}
	if repo.notifications[0].Icon != "store" || *repo.notifications[0].UserID != 10 {
		t.Errorf("expected store icon and userID 10, got icon=%s, user=%v", repo.notifications[0].Icon, repo.notifications[0].UserID)
	}

	// Calling again should not duplicate
	err = service.CreateWelcomeNotification(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.notifications) != 1 {
		t.Errorf("expected still 1 notification after duplicate call, got %d", len(repo.notifications))
	}
}
