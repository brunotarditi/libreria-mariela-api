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
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		response, err := c.service.Login(req)
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{
			"access_token": response.AccessToken,
			"info":         response.Info,
		})
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

func (c *UserController) FindAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := c.service.FindAll()
		if err != nil {
			ctx.JSON(500, gin.H{"error": "error al obtener usuarios"})
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

		ctx.JSON(200, usersResponse)
	}
}
