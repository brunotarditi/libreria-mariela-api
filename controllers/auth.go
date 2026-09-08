package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ctrl *AuthController) ExchangeToken(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El código OAuth es requerido"})
		return
	}

	// Obtener configuraciones del entorno
	peakAuthURL := os.Getenv("PEAK_AUTH_URL")
	if peakAuthURL == "" {
		peakAuthURL = "http://localhost:9009" // fallback por defecto local
	}
	clientSecret := os.Getenv("PEAK_AUTH_CLIENT_SECRET")
	clientID := "libreria-mariela" // Idealmente en ENV también

	// Canjear el código en Peak Auth
	resp, err := http.PostForm(peakAuthURL+"/oauth/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {req.Code},
		"grant_type":    {"authorization_code"},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error conectando con Peak Auth"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "El código OAuth es inválido o expiró"})
		return
	}

	// Reenviar la respuesta de Peak Auth (que incluye el access_token) a Angular
	var peakResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&peakResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo respuesta de Peak Auth"})
		return
	}

	c.JSON(http.StatusOK, peakResponse)
}
