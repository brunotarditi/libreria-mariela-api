package requests

import (
	"libreria/models"
)

type BudgetRequest struct {
	ClientName  string   `json:"client_name" binding:"required"`
	Description string   `json:"description"`
	Items       []string `json:"items" binding:"required"`
}

func (r BudgetRequest) ToModel() (models.Budget, error) {
	priceList := models.Budget{
		ClientName:  r.ClientName,
		Description: r.Description,
	}

	return priceList, nil
}

func (r BudgetRequest) UpdateModel(existing models.Budget) (models.Budget, error) {
	existing.ClientName = r.ClientName
	existing.Description = r.Description
	return existing, nil
}
