package controllers

import (
	"libreria/services"
	"net/http"
	"strconv"

	"github.com/brunotarditi/peak-auth/sdk/go"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	client              *peakauth.Client
	notificationService services.NotificationService
}

func NewAuthController(client *peakauth.Client, notificationService ...services.NotificationService) *AuthController {
	var notifService services.NotificationService
	if len(notificationService) > 0 {
		notifService = notificationService[0]
	}
	return &AuthController{client: client, notificationService: notifService}
}

func (ctrl *AuthController) ExchangeToken(c *gin.Context) {
	var req struct {
		Code         string `json:"code" binding:"required"`
		CodeVerifier string `json:"code_verifier" binding:"required"`
		RedirectURI  string `json:"redirect_uri" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan parámetros requeridos (code, code_verifier, redirect_uri)"})
		return
	}

	tokenResp, err := ctrl.client.ExchangeCode(c.Request.Context(), req.Code, req.CodeVerifier, req.RedirectURI)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Trigger automático: notificación de bienvenida al registrarse / login
	if ctrl.notificationService != nil && tokenResp.AccessToken != "" {
		if claims, err := ctrl.client.VerifyToken(tokenResp.AccessToken); err == nil && claims.Subject != "" {
			if uid64, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
				go func(uid uint) {
					_ = ctrl.notificationService.CreateWelcomeNotification(uid)
				}(uint(uid64))
			}
		}
	}

	c.JSON(http.StatusOK, tokenResp)
}
