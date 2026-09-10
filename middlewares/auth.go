package middlewares

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PeakAuthClaims representa el payload de los tokens JWT emitidos por Peak Auth.
type PeakAuthClaims struct {
	Username    string   `json:"username"`
	AppID       string   `json:"app_id"`
	Roles       []string `json:"roles"`
	MfaVerified bool     `json:"mfa_verified"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

// LoadRSAPublicKey carga y parsea la clave pública RSA desde variables de entorno o archivo PEM.
func LoadRSAPublicKey() (*rsa.PublicKey, error) {
	var pemBytes []byte

	// 1. Intentar leer clave pública directa desde variable de entorno PEAK_AUTH_PUBLIC_KEY
	if envPEM := os.Getenv("PEAK_AUTH_PUBLIC_KEY"); envPEM != "" {
		envPEM = strings.ReplaceAll(envPEM, "\\n", "\n")
		envPEM = strings.Trim(envPEM, "\"")
		pemBytes = []byte(envPEM)
	} else {
		// 2. Intentar leer desde la ruta especificada en PEAK_AUTH_PUBLIC_KEY_PATH o archivo local por defecto
		path := os.Getenv("PEAK_AUTH_PUBLIC_KEY_PATH")
		if path == "" {
			path = "jwt_public.pem"
		}

		if data, err := os.ReadFile(path); err == nil {
			pemBytes = data
		}
	}

	if len(pemBytes) == 0 {
		return nil, fmt.Errorf("no se encontró la clave pública de Peak Auth (configure PEAK_AUTH_PUBLIC_KEY o PEAK_AUTH_PUBLIC_KEY_PATH)")
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando la clave pública RSA desde PEM: %w", err)
	}

	return pubKey, nil
}

// AuthMiddleware valida el token JWT emitido por Peak Auth e inyecta la información del usuario en el contexto.
func AuthMiddleware(pubKey *rsa.PublicKey, expectedAppID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if pubKey == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Clave pública de autenticación no configurada en el servidor"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cabecera de autorización no provista"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Debe ser: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims := &PeakAuthClaims{}

		opts := []jwt.ParserOption{
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithLeeway(30 * time.Second),
		}

		// Validar que la audiencia coincida con la aplicación
		if expectedAppID != "" {
			opts = append(opts, jwt.WithAudience(expectedAppID))
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return pubKey, nil
		}, opts...)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// Inyectar claims en el contexto de Gin
		if claims.Subject != "" {
			if uid, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
				c.Set("user_id", uint(uid))
			}
		}
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("app_id", claims.AppID)

		c.Next()
	}
}

// RequireRole asegura que el usuario autenticado cuente con un rol específico.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesVal, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No se encontraron roles asignados para esta operación"})
			c.Abort()
			return
		}

		roles, ok := rolesVal.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Formato de roles inválido"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if strings.EqualFold(r, role) || strings.EqualFold(r, "ROOT") {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Acceso denegado: se requiere el rol '%s'", role)})
		c.Abort()
	}
}

// RequireAnyRole asegura que el usuario autenticado cuente con al menos uno de los roles solicitados.
func RequireAnyRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesVal, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No se encontraron roles asignados para esta operación"})
			c.Abort()
			return
		}

		roles, ok := rolesVal.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Formato de roles inválido"})
			c.Abort()
			return
		}

		for _, userRole := range roles {
			if strings.EqualFold(userRole, "ROOT") {
				c.Next()
				return
			}
			for _, allowed := range allowedRoles {
				if strings.EqualFold(userRole, allowed) {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: privilegios insuficientes"})
		c.Abort()
	}
}
