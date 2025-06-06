package requests

import "libreria/models"

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
