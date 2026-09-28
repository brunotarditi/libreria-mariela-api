package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/brunotarditi/peak-auth/sdk/go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetUserIDAndUserNameFromContext(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Sin claims ni contexto
	if uid := GetUserIDFromContext(c); uid != nil {
		t.Errorf("esperado nil, obtenido %v", uid)
	}
	if name := GetUserNameFromContext(c); name != "" {
		t.Errorf("esperado string vacío, obtenido %s", name)
	}

	// Con claims de Peak Auth
	claims := &peakauth.Claims{
		Username: "admin@libreria.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "15",
		},
	}
	c.Set("claims", claims)

	uid := GetUserIDFromContext(c)
	if uid == nil || *uid != 15 {
		t.Errorf("esperado userID 15, obtenido %v", uid)
	}

	name := GetUserNameFromContext(c)
	if name != "admin@libreria.com" {
		t.Errorf("esperado username admin@libreria.com, obtenido %s", name)
	}

	// Con subject pero sin username
	claimsNoName := &peakauth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "99",
		},
	}
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	c2.Set("claims", claimsNoName)

	if uid2 := GetUserIDFromContext(c2); uid2 == nil || *uid2 != 99 {
		t.Errorf("esperado userID 99, obtenido %v", uid2)
	}
	if name2 := GetUserNameFromContext(c2); name2 != "Usuario #99" {
		t.Errorf("esperado Usuario #99, obtenido %s", name2)
	}

	// Fallback gin context user_id y username
	c3, _ := gin.CreateTestContext(httptest.NewRecorder())
	c3.Set("user_id", uint(7))
	c3.Set("username", "empleado_local")

	if uid3 := GetUserIDFromContext(c3); uid3 == nil || *uid3 != 7 {
		t.Errorf("esperado userID 7, obtenido %v", uid3)
	}
	if name3 := GetUserNameFromContext(c3); name3 != "empleado_local" {
		t.Errorf("esperado username empleado_local, obtenido %s", name3)
	}
}
