package controllers

import (
	"fmt"
	"libreria/requests"
	"libreria/responses"
	"libreria/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserControllerController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := c.service.Login(req)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token": response.AccessToken,
			"info":         response.Info,
		})
	}
}

func (c *UserController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.RegisterRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := c.service.Register(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("el usuario %s ha sido registrado con éxito", user.Username)})
	}
}

func (c *UserController) AssignRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		val, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			return
		}

		currentUserID, ok := val.(uint)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar el ID del usuario"})
			return
		}

		userIDParam := ctx.Param("id")
		userID, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var request requests.Role
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.service.AssignRole(request, userID, currentUserID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "rol asignado con éxito"})
	}
}

func (c *UserController) FindAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := c.service.FindAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener usuarios"})
			return
		}

		if len(users) == 0 {
			ctx.JSON(http.StatusOK, []responses.UserResponse{})
			return
		}

		usersResponse := make([]responses.UserResponse, 0, len(users))
		for _, user := range users {
			userRes := responses.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				IsActive:  user.IsActive,
				LastLogin: user.LastLogin,
				CreatedAt: user.CreatedAt,
				Roles:     make([]string, 0, len(user.Roles)),
			}
			for _, role := range user.Roles {
				userRes.Roles = append(userRes.Roles, string(role.Name))
			}
			usersResponse = append(usersResponse, userRes)
		}

		ctx.JSON(http.StatusOK, usersResponse)
	}
}

func (c *UserController) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Falta el token"})
			return
		}

		err := c.service.VerifyEmail(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Correo verificado con éxito"})
	}
}

func (c *UserController) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Falta el email"})
			return
		}

		user, err := c.service.FindVerifiedUser(email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		allowed, err := c.service.CanRequestPasswordReset(user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando últimos resets"})
			return
		}
		if !allowed {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Ya solicitaste un restablecimiento recientemente, esperá unos minutos"})
			return
		}

		err = c.service.SendResetEmail(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Si el correo está registrado, vas a recibir instrucciones",
		})
	}
}

func (c *UserController) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.ResetPasswordRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		err := c.service.ResetPassword(req.Token, req.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
	}
}
