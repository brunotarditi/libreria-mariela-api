package repositories

import (
	"libreria/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (models.User, error)
	FindById(ID uint64) (models.User, error)
	Update(user *models.User) error
	UpdateColumn(column string, value interface{}, id uint) error
	VerifyRolExist(userID uint, roleID uint) int64
	CreateUserRole(userRole *models.UserRole) error
	RolesByUser(userID uint) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepositoryRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindById(id uint64) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateColumn(column string, value interface{}, id uint) error {
	return r.db.Model(&models.User{}).Where("id = ? AND is_active = ?", id, true).Update(column, value).Error
}

func (r *userRepository) VerifyRolExist(userID uint, roleID uint) int64 {
	var count int64
	r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count)

	return count
}

func (r *userRepository) CreateUserRole(userRole *models.UserRole) error {
	return r.db.Create(&userRole).Error
}

func (r *userRepository) RolesByUser(userID uint) (models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	return user, err
}
