package models

import (
	"gorm.io/gorm"
)

type BudgetItem struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type Budget struct {
	gorm.Model
	ClientName  string       `gorm:"type:varchar(100)"`
	Description string       `gorm:"type:text"`
	Items       []BudgetItem `gorm:"type:jsonb"`
}
