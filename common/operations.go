package common

import "gorm.io/gorm"

type Operations[T any] interface {
	FindAll() ([]T, error)
	FindByID(id string) (T, error)
	Create(model T) error
	Update(model T) error
	Delete(id string) error
}

type GormOperations[T any] struct {
	db *gorm.DB
}

func NewGormOperations[T any](db *gorm.DB) *GormOperations[T] {
	return &GormOperations[T]{db: db}
}

func (ops *GormOperations[T]) FindAll() ([]T, error) {
	var models []T
	err := ops.db.Find(&models).Error
	return models, err
}

func (ops *GormOperations[T]) FindByID(id string) (T, error) {
	var model T
	err := ops.db.First(&model, id).Error
	return model, err
}

func (ops *GormOperations[T]) Create(model T) error {
	return ops.db.Create(&model).Error
}

func (ops *GormOperations[T]) Update(model T) error {
	return ops.db.Save(&model).Error
}

func (ops *GormOperations[T]) Delete(id string) error {
	var model T
	return ops.db.Delete(&model, id).Error
}
