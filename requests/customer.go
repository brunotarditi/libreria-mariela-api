package requests

import "libreria/models"

type CustomerRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=65"`
	ContactInfo string `json:"contact_info" binding:"required,min=1,max=255"`
}

func (r CustomerRequest) ToModel() (models.Customer, error) {
	return models.Customer{
		Name:        r.Name,
		ContactInfo: r.ContactInfo,
	}, nil
}

func (r CustomerRequest) UpdateModel(existing models.Customer) (models.Customer, error) {
	existing.Name = r.Name
	existing.ContactInfo = r.ContactInfo
	return existing, nil
}
