package middlewares

import (
	"libreria/models"
	"libreria/utils"
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

		// Recuperar el ID de usuario y nombre/email desde el contexto o claims
		userID := utils.GetUserIDFromContext(c)
		userName := utils.GetUserNameFromContext(c)

		// Continuar con el procesamiento
		c.Next()

		audit := models.AuditLog{
			UserID:    userID,
			UserName:  userName,
			Route:     route,
			Method:    method,
			IP:        ip,
			RequestAt: start,
		}

		// Guardar en BD si la instancia de base de datos está disponible
		if db != nil {
			go db.Create(&audit)
		}
	}
}
