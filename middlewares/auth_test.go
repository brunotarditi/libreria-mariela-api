package middlewares

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func generateTestKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generando par de claves RSA de prueba: %v", err)
	}
	return privKey, &privKey.PublicKey
}

func createTestToken(t *testing.T, privKey *rsa.PrivateKey, appID string, userID string, roles []string, exp time.Duration) string {
	t.Helper()
	claims := PeakAuthClaims{
		Username:    "testuser@libreria.com",
		AppID:       appID,
		Roles:       roles,
		MfaVerified: true,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "peak-auth",
			Audience:  jwt.ClaimStrings{appID},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(privKey)
	if err != nil {
		t.Fatalf("error firmando token de prueba: %v", err)
	}
	return tokenStr
}

func TestAuthMiddleware(t *testing.T) {
	privKey, pubKey := generateTestKeyPair(t)
	expectedApp := "libreria-mariela"

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware(pubKey, expectedApp))
		r.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			roles, _ := c.Get("roles")
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"roles":    roles,
				"username": username,
			})
		})
		return r
	}

	t.Run("Sin cabecera Authorization", func(t *testing.T) {
		r := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("esperado status %d, obtenido %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Formato de cabecera inválido", func(t *testing.T) {
		r := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "TokenInvalidoSinBearer")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("esperado status %d, obtenido %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Token expirado", func(t *testing.T) {
		r := setupRouter()
		expiredToken := createTestToken(t, privKey, expectedApp, "10", []string{"USER"}, -10*time.Minute)
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("esperado status %d, obtenido %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Token con app_id distinta", func(t *testing.T) {
		r := setupRouter()
		otherAppToken := createTestToken(t, privKey, "otra-app", "10", []string{"USER"}, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+otherAppToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("esperado status %d por audiencia errónea, obtenido %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Token válido", func(t *testing.T) {
		r := setupRouter()
		validToken := createTestToken(t, privKey, expectedApp, "25", []string{"ADMIN", "CAJERO"}, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("esperado status %d, obtenido %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestRequireRole(t *testing.T) {
	privKey, pubKey := generateTestKeyPair(t)
	expectedApp := "libreria-mariela"

	setupRouterWithRole := func(requiredRole string) *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware(pubKey, expectedApp))
		r.GET("/admin-only", RequireRole(requiredRole), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		return r
	}

	t.Run("Acceso denegado si no tiene el rol", func(t *testing.T) {
		r := setupRouterWithRole("ADMIN")
		userToken := createTestToken(t, privKey, expectedApp, "5", []string{"CAJERO"}, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("esperado status %d, obtenido %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Acceso concedido si tiene el rol", func(t *testing.T) {
		r := setupRouterWithRole("ADMIN")
		adminToken := createTestToken(t, privKey, expectedApp, "1", []string{"ADMIN"}, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("esperado status %d, obtenido %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Acceso concedido para rol ROOT automáticamente", func(t *testing.T) {
		r := setupRouterWithRole("ADMIN")
		rootToken := createTestToken(t, privKey, expectedApp, "1", []string{"ROOT"}, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		req.Header.Set("Authorization", "Bearer "+rootToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("esperado status %d para ROOT, obtenido %d", http.StatusOK, w.Code)
		}
	})
}
