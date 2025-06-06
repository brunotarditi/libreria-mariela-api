package middlewares

import (
	"libreria/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var roleHierarchy = map[string]int{
	"READ":  1,
	"WRITE": 2,
	"ADMIN": 3,
}

func RoleMiddleware(db *gorm.DB, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idAny, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		userID := idAny.(uint)

		var userRoles []models.Role
		db.Model(&models.Role{}).
			Joins("JOIN user_roles ON user_roles.role_id = roles.id").
			Where("user_roles.user_id = ?", userID).
			Find(&userRoles)

		for _, role := range userRoles {
			if roleHierarchy[role.Name] >= roleHierarchy[requiredRole] {
				c.Next()
				return
			}
		}

		c.Next()
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool)
	for _, role := range roles {
		roleSet[strings.ToUpper(role)] = true
	}

	return func(c *gin.Context) {
		val, exists := c.Get("user_roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No se encontraron roles en el contexto"})
			return
		}

		userRoles := val.([]string)
		for _, role := range userRoles {
			if roleSet[strings.ToUpper(role)] {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No tienes el rol requerido para tal acción"})
	}
}
