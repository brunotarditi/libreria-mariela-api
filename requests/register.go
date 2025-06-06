package requests

import (
	"libreria/models"
	"libreria/utils"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r RegisterRequest) ToUser() (models.User, error) {
	hashedPassword, err := utils.HashPassword(r.Password)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Username: r.Username,
		Email:    r.Email,
		Password: hashedPassword,
		IsActive: true,
	}, nil
}
