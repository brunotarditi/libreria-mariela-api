package common

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type QueryOptions struct {
	Search    string
	Sort      string
	Direction string
}

type Operations[T any] interface {
	FindAll() ([]T, error)
	FindByID(id string) (T, error)
	Paginated(offset, size int, options QueryOptions) ([]T, error)
	Count(options QueryOptions) (int64, error)
	Create(model T) error
	CreateMany(model []T) error
	Update(model T) error
	Delete(id string) error
	Pluck(field string) ([]string, error)
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

func (ops *GormOperations[T]) Paginated(offset, size int, options QueryOptions) ([]T, error) {
	var models []T
	query := ops.buildQuery(options)
	err := query.Offset(offset).Limit(size).Find(&models).Error
	return models, err
}

func (ops *GormOperations[T]) Count(options QueryOptions) (int64, error) {
	var total int64
	query := ops.buildQuery(options)
	err := query.Count(&total).Error
	return total, err
}

func (ops *GormOperations[T]) Create(model T) error {
	return ops.db.Create(&model).Error
}

func (ops *GormOperations[T]) CreateMany(model []T) error {
	return ops.db.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(model, 100).Error
	})
}

func (ops *GormOperations[T]) Update(model T) error {
	return ops.db.Save(&model).Error
}

func (ops *GormOperations[T]) Delete(id string) error {
	var model T
	return ops.db.Delete(&model, id).Error
}

func (ops *GormOperations[T]) Pluck(field string) ([]string, error) {
	var result []string
	err := ops.db.Model(new(T)).Pluck(field, &result).Error
	return result, err
}

func (ops *GormOperations[T]) buildQuery(options QueryOptions) *gorm.DB {
	query := ops.db.Model(new(T))

	if options.Search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(options.Search)+"%")
	}

	if options.Sort == "name" && (options.Direction == "asc" || options.Direction == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", options.Sort, options.Direction))
	}

	return query
}
