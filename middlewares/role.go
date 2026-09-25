package middlewares

import (
	"net/http"

	peakauthgin "github.com/brunotarditi/peak-auth/sdk/go/gin"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware protege las rutas basándose en el método HTTP y los roles del JWT.
func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := peakauthgin.ClaimsFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		roles := claims.Roles
		hasRole := func(r string) bool {
			for _, role := range roles {
				if role == r || role == "ROOT" {
					return true
				}
			}
			return false
		}

		// GET requiere READ, WRITE o ADMIN
		if c.Request.Method == http.MethodGet {
			if !hasRole("READ") && !hasRole("WRITE") && !hasRole("ADMIN") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permisos de lectura (READ) requeridos"})
				return
			}
		}

		// POST/PUT/PATCH requiere WRITE o ADMIN
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			if !hasRole("WRITE") && !hasRole("ADMIN") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permisos de escritura (WRITE) requeridos"})
				return
			}
		}

		// DELETE requiere ADMIN (se permite WRITE para gestionar notificaciones propias)
		if c.Request.Method == http.MethodDelete {
			if c.FullPath() == "/api/v1/notifications" || c.FullPath() == "/api/v1/notifications/:id" {
				if !hasRole("WRITE") && !hasRole("ADMIN") {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permisos de escritura (WRITE) requeridos"})
					return
				}
			} else if !hasRole("ADMIN") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permisos de administrador (ADMIN) requeridos"})
				return
			}
		}

		c.Next()
	}
}
