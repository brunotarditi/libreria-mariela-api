package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	Name            string            `gorm:"type:varchar(65);not null" json:"name"`
	ContactInfo     string            `gorm:"type:varchar(255)" json:"contact_info"`
	PurchaseHistory []PurchaseHistory `gorm:"foreignKey:SupplierID" json:"-"`
}
