package controllers

import (
	"fmt"
	"libreria/requests"
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
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		token, err := c.service.Login(req)
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"token": token})
	}
}

func (c *UserController) UserRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.RegisterRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user, err := c.service.Register(request)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": fmt.Sprintf("el usuario %s ha sido registrado con éxito", user.Username)})
	}
}

func (c *UserController) AssignRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userIDParam := ctx.Param("id")
		userID, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var request requests.Role
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := c.service.AssignRole(request, userID); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "rol asignado con éxito"})
	}
}
