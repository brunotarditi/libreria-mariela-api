package requests

import "libreria/models"

type SupplierRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=65"`
	ContactInfo string `json:"contact_info" binding:"required,min=1,max=255"`
}

func (r SupplierRequest) ToModel() (models.Supplier, error) {
	return models.Supplier{
		Name:        r.Name,
		ContactInfo: r.ContactInfo,
	}, nil
}

func (r SupplierRequest) UpdateModel(existing models.Supplier) (models.Supplier, error) {
	existing.Name = r.Name
	existing.ContactInfo = r.ContactInfo
	return existing, nil
}
