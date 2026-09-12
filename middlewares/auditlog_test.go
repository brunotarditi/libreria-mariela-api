package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunotarditi/peak-auth/sdk/go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuditMiddlewareWithPeakAuthClaims(t *testing.T) {
	r := gin.New()

	// Simulamos el middleware del SDK inyectando los claims antes de AuditMiddleware
	r.Use(func(c *gin.Context) {
		claims := &peakauth.Claims{
			Username: "empleado@libreria.com",
			AppID:    "libreria-mariela",
			Roles:    []string{"CAJERO"},
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "42",
			},
		}
		c.Set("claims", claims)
		c.Set("user", claims)
		c.Next()
	})

	r.Use(AuditMiddleware(nil))

	r.GET("/test-audit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test-audit", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado status %d, obtenido %d", http.StatusOK, w.Code)
	}
}
