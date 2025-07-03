package requests

import (
	"errors"
	"libreria/models"

	"gorm.io/gorm"
)

type CategoryRequestArray []CategoryRequest
type CategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=65"`
}

func (r CategoryRequest) ToModel() (models.Category, error) {
	return models.Category{Name: r.Name}, nil
}

func (r CategoryRequest) UpdateModel(existing models.Category) (models.Category, error) {
	existing.Name = r.Name
	return existing, nil
}

func (r CategoryRequestArray) ToArrayModel() ([]models.Category, error) {
	categories := make([]models.Category, 0, len(r))
	for _, req := range r {
		category, err := req.ToModel()
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r CategoryRequest) Validate(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Category{}).Where("name = ?", r.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("el nombre de la categoría ya existe")
	}
	return nil
}

func (r CategoryRequest) ValidateUpdate(db *gorm.DB, existing models.Category) error {
	if r.Name == existing.Name {
		return nil
	}

	var count int64
	if err := db.Model(&models.Category{}).
		Where("name = ?", r.Name).
		Where("id != ?", existing.ID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe otra categoría con ese nombre")
	}
	return nil
}

func (r CategoryRequestArray) Validate(db *gorm.DB) error {
	for _, item := range r {
		if v, ok := any(item).(ValidateRequest); ok {
			if err := v.Validate(db); err != nil {
				return err
			}
		}
	}
	return nil
}
