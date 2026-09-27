package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunotarditi/peak-auth/sdk/go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupTestRouterWithRoles(roles []string) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if roles != nil {
			claims := &peakauth.Claims{
				Username: "user@libreria.com",
				AppID:    "libreria-mariela",
				Roles:    roles,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "42",
				},
			}
			c.Set("claims", claims)
		}
		c.Next()
	})
	r.Use(RoleMiddleware())

	r.GET("/api/v1/products", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/api/v1/products", func(c *gin.Context) { c.Status(http.StatusCreated) })
	r.DELETE("/api/v1/products/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.DELETE("/api/v1/notifications/:id", func(c *gin.Context) { c.Status(http.StatusOK) })

	return r
}

func TestRoleMiddleware_OwnerAccess(t *testing.T) {
	r := setupTestRouterWithRoles([]string{"OWNER"})

	// OWNER can access GET, POST, DELETE
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for OWNER GET, got %d", w.Code)
	}

	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for OWNER DELETE, got %d", w.Code)
	}
}

func TestRoleMiddleware_AdminAccess(t *testing.T) {
	r := setupTestRouterWithRoles([]string{"ADMIN"})

	// ADMIN can access DELETE
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for ADMIN DELETE, got %d", w.Code)
	}
}

func TestRoleMiddleware_UserAccess(t *testing.T) {
	r := setupTestRouterWithRoles([]string{"USER"})

	// USER can access GET
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for USER GET, got %d", w.Code)
	}

	// USER can access POST
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/products", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 for USER POST, got %d", w.Code)
	}

	// USER can delete own notification
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/notifications/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for USER DELETE notification, got %d", w.Code)
	}

	// USER cannot delete a product (requires ADMIN or OWNER)
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for USER DELETE product, got %d", w.Code)
	}
}

func TestRoleMiddleware_UnauthorizedWithoutClaims(t *testing.T) {
	r := setupTestRouterWithRoles(nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing claims, got %d", w.Code)
	}
}
