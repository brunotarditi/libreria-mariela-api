package controllers

import (
	"libreria/services"
	"libreria/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service services.NotificationService
}

func NewNotificationController(service services.NotificationService) *NotificationController {
	return &NotificationController{service: service}
}

func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	unreadOnly := c.Query("unread") == "true"
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	userID := utils.GetUserIDFromContext(c)

	notifications, unreadCount, err := ctrl.service.GetNotifications(userID, unreadOnly, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         notifications,
		"total":        len(notifications),
		"unread_count": unreadCount,
	})
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	if err := ctrl.service.MarkAsRead(uint(id64), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notificación marcada como leída",
	})
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	count, err := ctrl.service.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Todas las notificaciones marcadas como leídas",
		"updated_count": count,
	})
}

func (ctrl *NotificationController) Delete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	if err := ctrl.service.Delete(uint(id64), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notificación eliminada con éxito",
	})
}

func (ctrl *NotificationController) ClearAll(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	count, err := ctrl.service.ClearAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Todas las notificaciones fueron eliminadas",
		"deleted_count": count,
	})
}
