package utils

import (
	"strconv"

	peakauthgin "github.com/brunotarditi/peak-auth/sdk/go/gin"
	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext retrieves the user ID from the Peak Auth claims or Gin context.
func GetUserIDFromContext(c *gin.Context) *uint {
	if claims, ok := peakauthgin.ClaimsFromContext(c); ok && claims.Subject != "" {
		if uid64, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
			uid := uint(uid64)
			return &uid
		}
	} else if id, exists := c.Get("user_id"); exists {
		if uid, ok := id.(uint); ok {
			return &uid
		}
	}
	return nil
}

// GetUserNameFromContext retrieves the username or subject representation from Peak Auth claims or Gin context.
func GetUserNameFromContext(c *gin.Context) string {
	if claims, ok := peakauthgin.ClaimsFromContext(c); ok {
		if claims.Username != "" {
			return claims.Username
		}
		if claims.Subject != "" {
			return "Usuario #" + claims.Subject
		}
	}
	if name, exists := c.Get("username"); exists {
		if uName, ok := name.(string); ok && uName != "" {
			return uName
		}
	}
	if name, exists := c.Get("user_name"); exists {
		if uName, ok := name.(string); ok && uName != "" {
			return uName
		}
	}
	return ""
}

