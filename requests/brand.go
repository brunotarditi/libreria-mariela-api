package requests

import "libreria/models"

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
