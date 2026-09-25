package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	FindAll(userID *uint, unreadOnly bool, limit int) ([]models.Notification, error)
	CountUnread(userID *uint) (int64, error)
	FindByID(id uint) (models.Notification, error)
	MarkAsRead(id uint, userID *uint) error
	MarkAllAsRead(userID *uint) (int64, error)
	Delete(id uint, userID *uint) error
	DeleteAll(userID *uint) (int64, error)
	Create(notification *models.Notification) error
	ExistsUnreadByRoute(route string) (bool, error)
	ExistsWelcome(userID uint) (bool, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) FindAll(userID *uint, unreadOnly bool, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := r.db.Model(&models.Notification{})

	if userID != nil {
		query = query.Where("user_id IS NULL OR user_id = ?", *userID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	query = query.Order("created_at DESC")

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query = query.Limit(limit)

	err := query.Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountUnread(userID *uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Notification{}).Where("is_read = ?", false)

	if userID != nil {
		query = query.Where("user_id IS NULL OR user_id = ?", *userID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *notificationRepository) FindByID(id uint) (models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	return notification, err
}

func (r *notificationRepository) MarkAsRead(id uint, userID *uint) error {
	query := r.db.Model(&models.Notification{}).Where("id = ?", id)
	if userID != nil {
		query = query.Where("user_id IS NULL OR user_id = ?", *userID)
	}
	return query.Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID *uint) (int64, error) {
	query := r.db.Model(&models.Notification{}).Where("is_read = ?", false)
	if userID != nil {
		query = query.Where("user_id IS NULL OR user_id = ?", *userID)
	}
	res := query.Update("is_read", true)
	return res.RowsAffected, res.Error
}

func (r *notificationRepository) Delete(id uint, userID *uint) error {
	query := r.db.Where("id = ?", id)
	if userID != nil {
		query = query.Where("user_id IS NULL OR user_id = ?", *userID)
	}
	return query.Delete(&models.Notification{}).Error
}

func (r *notificationRepository) DeleteAll(userID *uint) (int64, error) {
	query := r.db.Model(&models.Notification{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	res := query.Delete(&models.Notification{})
	return res.RowsAffected, res.Error
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) ExistsUnreadByRoute(route string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("route = ? AND is_read = ?", route, false).
		Count(&count).Error
	return count > 0, err
}

func (r *notificationRepository) ExistsWelcome(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND icon = 'store'", userID).
		Count(&count).Error
	return count > 0, err
}
