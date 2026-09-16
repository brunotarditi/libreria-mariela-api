package controllers

import (
	"net/http"

	"github.com/brunotarditi/peak-auth/sdk/go"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	client *peakauth.Client
}

func NewAuthController(client *peakauth.Client) *AuthController {
	return &AuthController{client: client}
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

	c.JSON(http.StatusOK, tokenResp)
}
