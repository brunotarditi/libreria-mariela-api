package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"libreria/models"
	"libreria/repositories"
	"libreria/requests"
	"libreria/responses"
	"libreria/security"
	"libreria/utils"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Register(req requests.RegisterRequest) (models.User, error)
	Login(req requests.LoginRequest) (responses.TokenResponse, error)
	AssignRole(req requests.Role, id uint64, currentUserID uint) error
	FindRolesByUser(userID uint) (models.User, error)
	FindAll() ([]models.User, error)
	VerifyEmail(token string) error
	ResetPassword(token, newPassword string) error
	FindVerifiedUser(email string) (*models.User, error)
	CanRequestPasswordReset(userID uint) (bool, error)
	SendResetEmail(user *models.User) error
}

type userService struct {
	db                    *gorm.DB
	userRepo              repositories.UserRepository
	roleRepo              repositories.RoleRepository
	tokenManager          *security.PasetoManager
	emailVerificationRepo repositories.EmailVerificationRepository
	passwordResetRepo     repositories.PasswordResetRepository
}

func NewUserService(db *gorm.DB, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, tokenManager *security.PasetoManager, emailVerificationRepo repositories.EmailVerificationRepository, passwordResetRepo repositories.PasswordResetRepository) UserService {
	return &userService{db: db, userRepo: userRepo, roleRepo: roleRepo, tokenManager: tokenManager, emailVerificationRepo: emailVerificationRepo, passwordResetRepo: passwordResetRepo}
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

	if !user.IsVerified {
		return responses.TokenResponse{}, fmt.Errorf("usuario no verificado")
	}

	if !user.IsActive {
		return responses.TokenResponse{}, fmt.Errorf("usuario está desactivado")
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

	if err := utils.ValidatePassword(req.Password); err != nil {
		return models.User{}, err
	}

	user, err := req.ToUser()
	if err != nil {
		return models.User{}, err
	}

	if exists, _ := s.userRepo.ExistsByEmail(req.Email); exists {
		return models.User{}, fmt.Errorf("el email ya está registrado")
	}

	if exists, _ := s.userRepo.ExistsByUsername(req.Username); exists {
		return models.User{}, fmt.Errorf("el username ya está en uso")
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

	plainToken, tokenHash, err := utils.GenerateToken(32)
	if err != nil {
		return models.User{}, err
	}

	verification := models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.emailVerificationRepo.CreateEmailVerification(&verification); err != nil {
		return models.User{}, err
	}

	tx.Commit()
	url := os.Getenv("VERIFY_URL")
	verifyURL := fmt.Sprintf(url, plainToken)

	htmlBody, err := utils.RenderVerificationEmail("templates/verify-email.html", map[string]string{
		"Username":  user.Username,
		"VerifyURL": verifyURL,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("error al renderizar email: %w", err)
	}

	if err := sendVerificationEmail("Verificá tu email", user.Email, htmlBody); err != nil {
		return models.User{}, fmt.Errorf("error enviando email: %v", err)
	}
	return user, nil
}

func (s *userService) AssignRole(req requests.Role, userID uint64, currentUserID uint) error {

	if userID == uint64(currentUserID) {
		return fmt.Errorf("no puedes asignarte roles a ti mismo")
	}

	if strings.ToUpper(req.RoleName) == "ROOT" {
		count, err := s.userRepo.CountUsersWithRole("ROOT")
		if err != nil {
			return err
		}
		if count >= 1 {
			return fmt.Errorf("ya existe un usuario con rol ROOT")
		}
	}

	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

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

func (s *userService) FindAll() ([]models.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %v", err)
	}
	return users, nil
}

func (s *userService) VerifyEmail(token string) error {

	verification, err := s.emailVerificationRepo.FindEmailVerification(token)
	if err != nil {
		return fmt.Errorf("token inválido o expirado")
	}

	tx := s.db.Begin()
	now := time.Now()
	if err := s.emailVerificationRepo.UpdateUsedAt(verification, now); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.userRepo.UpdateColumn("is_verified", true, verification.UserID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *userService) FindVerifiedUser(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("usuario no verificado")
	}
	return &user, nil
}

func (s *userService) CanRequestPasswordReset(userID uint) (bool, error) {
	lastReset, err := s.passwordResetRepo.CheckLastTimeTokenReset(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return time.Since(lastReset) >= 15*time.Minute, nil
}

func (s *userService) SendResetEmail(user *models.User) error {
	plainToken, tokenHash, err := utils.GenerateToken(32)
	if err != nil {
		return err
	}

	reset := &models.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := s.passwordResetRepo.CreatePasswordReset(reset); err != nil {
		return err
	}

	verifyURL := fmt.Sprintf(os.Getenv("RESET_PASSWORD"), plainToken)
	htmlBody, err := utils.RenderVerificationEmail("templates/reset-password.html", map[string]string{
		"Username":  user.Username,
		"VerifyURL": verifyURL,
	})
	if err != nil {
		return fmt.Errorf("error al renderizar email: %w", err)
	}

	if err := sendVerificationEmail("Restablecer contraseña", user.Email, htmlBody); err != nil {
		return fmt.Errorf("error enviando email: %v", err)
	}
	return nil
}

func (s *userService) ResetPassword(token, newPassword string) error {
	reset, err := s.passwordResetRepo.FindValidPasswordReset(token)
	if err != nil {
		return fmt.Errorf("token inválido o expirado")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error al hashear contraseña: %w", err)
	}

	if err := s.passwordResetRepo.UpdatePassword(reset.UserID, hashed); err != nil {
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	now := time.Now()
	if err := s.passwordResetRepo.MarkPasswordResetUsed(reset.ID, now); err != nil {
		return fmt.Errorf("error al actualizar estado del token: %w", err)
	}

	return nil
}
