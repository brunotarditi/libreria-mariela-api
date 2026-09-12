package middlewares

import (
	"libreria/models"
	"strconv"
	"time"

	peakauthgin "github.com/brunotarditi/peak-auth/sdk/go/gin"
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

		// Recuperar el ID de usuario desde los claims del SDK de Peak Auth
		var userID *uint = nil
		if claims, ok := peakauthgin.ClaimsFromContext(c); ok && claims.Subject != "" {
			if uid64, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
				uid := uint(uid64)
				userID = &uid
			}
		} else if id, exists := c.Get("user_id"); exists {
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

		// Guardar en BD si la instancia de base de datos está disponible
		if db != nil {
			go db.Create(&audit)
		}
	}
}
