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
