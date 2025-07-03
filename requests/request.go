package requests

import "gorm.io/gorm"

type MapperRequest[T any] interface {
	ToModel() (T, error)
	UpdateModel(T) (T, error)
}

type ValidateRequest interface {
	Validate(*gorm.DB) error
}

type ValidateOnUpdate[T any] interface {
	ValidateUpdate(*gorm.DB, T) error
}

type MapperArrayRequest[T any] interface {
	ToArrayModel() ([]T, error)
}
