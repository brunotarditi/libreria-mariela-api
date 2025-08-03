package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	Create(user *models.User) error
	FindByEmail(email string) (models.User, error)
	FindById(ID uint64) (models.User, error)
	Update(user *models.User) error
	UpdateColumn(column string, value interface{}, id uint) error
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	VerifyRolExist(userID uint, roleID uint) int64
	CreateUserRole(userRole *models.UserRole) error
	CountUsersWithRole(roleName string) (int64, error)
	RolesByUser(userID uint) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepositoryRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Roles").Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindById(id uint64) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateColumn(column string, value interface{}, id uint) error {
	return r.db.Model(&models.User{}).Where("id = ? AND is_active = ?", id, true).Update(column, value).Error
}

func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) VerifyRolExist(userID uint, roleID uint) int64 {
	var count int64
	r.db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	return count
}

func (r *userRepository) CreateUserRole(userRole *models.UserRole) error {
	return r.db.Create(&userRole).Error
}

func (r *userRepository) CountUsersWithRole(roleName string) (int64, error) {
	var count int64
	err := r.db.Table("user_roles").Joins("JOIN roles ON roles.id = user_roles.role_id").Where("roles.name = ?", roleName).Count(&count).Error

	return count, err
}

func (r *userRepository) RolesByUser(userID uint) (models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	return user, err
}
