package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"
	"libreria/responses"
	"libreria/security"
	"libreria/utils"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	Register(req requests.RegisterRequest) (models.User, error)
	Login(req requests.LoginRequest) (responses.TokenResponse, error)
	AssignRole(req requests.Role, id uint64) error
	FindRolesByUser(userID uint) (models.User, error)
}

type userService struct {
	db           *gorm.DB
	userRepo     repositories.UserRepository
	roleRepo     repositories.RoleRepository
	tokenManager *security.PasetoManager
}

func NewUserService(db *gorm.DB, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, tokenManager *security.PasetoManager) UserService {
	return &userService{db: db, userRepo: userRepo, roleRepo: roleRepo, tokenManager: tokenManager}
}

func (s *userService) Login(req requests.LoginRequest) (responses.TokenResponse, error) {
	var user models.User

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return responses.TokenResponse{}, fmt.Errorf("usuario no encontrado")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return responses.TokenResponse{}, fmt.Errorf("contraseña incorrecta")
	}

	roleNames := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleNames[i] = r.Name
	}

	rolesJSON, err := json.Marshal(roleNames)
	if err != nil {
		fmt.Println("Error al serializar:", err)
		return responses.TokenResponse{}, fmt.Errorf("error al serializar: %v", err)
	}

	info := base64.StdEncoding.EncodeToString(rolesJSON)
	token, err := s.tokenManager.GenerateToken(user.ID, user.Username, time.Hour*12)
	if err != nil {
		return responses.TokenResponse{}, fmt.Errorf("error generando token: %v", err)
	}

	user.LastLogin = time.Now()

	tx := s.db.Begin()
	if err := s.userRepo.UpdateColumn("last_login", user.LastLogin, user.ID); err != nil {
		tx.Rollback()
		return responses.TokenResponse{}, err
	}

	return responses.TokenResponse{
		AccessToken: token,
		Info:        info,
	}, nil
}

func (s *userService) Register(req requests.RegisterRequest) (models.User, error) {

	user, err := req.ToUser()
	if err != nil {
		return models.User{}, err
	}

	tx := s.db.Begin()

	if err := s.userRepo.Create(&user); err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	role, err := s.roleRepo.FindByRoleName("READ")
	if err != nil {
		return models.User{}, err
	}

	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}

	if err := s.userRepo.CreateUserRole(&userRole); err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	tx.Commit()
	return user, nil
}

func (s *userService) AssignRole(req requests.Role, userID uint64) error {

	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Verificamos que el rol exista
	role, err := s.roleRepo.FindByRoleName(req.RoleName)
	if err != nil {
		return fmt.Errorf("rol no encontrado")
	}

	if s.userRepo.VerifyRolExist(user.ID, role.ID) > 0 {
		return fmt.Errorf("ese rol ya está asignado")
	}

	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}

	tx := s.db.Begin()
	if err := s.userRepo.CreateUserRole(&userRole); err != nil {
		tx.Rollback()
		return fmt.Errorf("falló en la asignación del rol")
	}

	tx.Commit()
	return nil
}

func (s *userService) FindRolesByUser(userID uint) (models.User, error) {
	user, err := s.userRepo.RolesByUser(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("roles no encontrados")
	}
	return user, nil
}
