package middlewares

import (
	"fmt"
	"libreria/security"
	"libreria/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(manager *security.PasetoManager, userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token no provisto"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		jsonToken, err := manager.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		// extraer userID del Subject
		var userID uint
		_, err = fmt.Sscanf(jsonToken.Subject, "%d", &userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
			return
		}

		// agregar al contexto
		c.Set("user_id", userID)
		c.Set("username", jsonToken.Get("username"))

		user, err := userService.FindRolesByUser(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
			return
		}

		var roleNames []string
		for _, role := range user.Roles {
			roleNames = append(roleNames, role.Name)
		}

		c.Set("user_roles", roleNames)

		c.Next()
	}
}
