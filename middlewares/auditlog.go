package middlewares

import (
	"libreria/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Antes de procesar la request
		start := time.Now()
		ip := c.ClientIP()
		route := c.FullPath()
		method := c.Request.Method

		// Si después usás auth, podés recuperar el ID del usuario del contexto
		var userID *uint = nil
		if id, exists := c.Get("user_id"); exists {
			uid := id.(uint)
			userID = &uid
		}

		// Continuar con el procesamiento
		c.Next()

		audit := models.AuditLog{
			UserID:    userID,
			Route:     route,
			Method:    method,
			IP:        ip,
			RequestAt: start,
		}

		// Guardar en BD
		go db.Create(&audit) // lo hacemos en goroutine para no bloquear la request
	}
}
