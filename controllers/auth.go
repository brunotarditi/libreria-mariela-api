package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

	// Obtener y validar configuraciones del entorno
	clientSecret := os.Getenv("PEAK_AUTH_CLIENT_SECRET")
	if clientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuración OAuth incompleta en el servidor: PEAK_AUTH_CLIENT_SECRET no está definido"})
		return
	}

	peakAuthURL := os.Getenv("PEAK_AUTH_URL")
	if peakAuthURL == "" {
		peakAuthURL = "http://localhost:9009" // fallback por defecto local
	}

	clientID := os.Getenv("PEAK_AUTH_CLIENT_ID")
	if clientID == "" {
		clientID = "libreria-mariela"
	}

	// Preparar datos del formulario para el canje
	formData := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {req.Code},
		"grant_type":    {"authorization_code"},
	}

	reqUpstream, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		peakAuthURL+"/oauth/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar la petición hacia Peak Auth"})
		return
	}
	reqUpstream.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Cliente HTTP con timeout explícito
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(reqUpstream)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Error al conectar con Peak Auth: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= 500 {
			c.JSON(http.StatusBadGateway, gin.H{"error": "El servicio Peak Auth no está disponible temporalmente"})
			return
		}
		c.JSON(resp.StatusCode, gin.H{"error": "El código OAuth es inválido o expiró"})
		return
	}

	// Reenviar la respuesta de Peak Auth (que incluye el access_token) al cliente
	var peakResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&peakResponse); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Error leyendo respuesta de Peak Auth"})
		return
	}

	c.JSON(http.StatusOK, peakResponse)
}
