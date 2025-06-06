package requests

import "gorm.io/gorm"

type MapperRequest[T any] interface {
	ToModel() (T, error)
	UpdateModel(T) (T, error)
}

type ValidateRequest interface {
	Validate(*gorm.DB) error
}
