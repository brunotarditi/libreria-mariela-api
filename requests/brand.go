package requests

import (
	"errors"
	"libreria/models"

	"gorm.io/gorm"
)

type BrandRequestArray []BrandRequest
type BrandRequest struct {
	Name string `json:"name" binding:"required,min=1,max=65"`
}

func (r BrandRequest) ToModel() (models.Brand, error) {
	return models.Brand{Name: r.Name}, nil
}

func (r BrandRequest) UpdateModel(existing models.Brand) (models.Brand, error) {
	existing.Name = r.Name
	return existing, nil
}

func (r BrandRequestArray) ToArrayModel() ([]models.Brand, error) {
	brands := make([]models.Brand, 0, len(r))
	for _, req := range r {
		brand, err := req.ToModel()
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	return brands, nil
}

func (r BrandRequest) Validate(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Brand{}).Where("name = ?", r.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("el nombre de la marca ya existe")
	}
	return nil
}

func (r BrandRequest) ValidateUpdate(db *gorm.DB, existing models.Brand) error {
	if r.Name == existing.Name {
		return nil
	}

	var count int64
	if err := db.Model(&models.Brand{}).
		Where("name = ?", r.Name).
		Where("id != ?", existing.ID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe otra marca con ese nombre")
	}
	return nil
}

func (r BrandRequestArray) Validate(db *gorm.DB) error {
	for _, item := range r {
		if v, ok := any(item).(ValidateRequest); ok {
			if err := v.Validate(db); err != nil {
				return err
			}
		}
	}
	return nil
}
