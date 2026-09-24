package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	Name            string            `gorm:"type:varchar(65);not null" json:"name"`
	ContactInfo     string            `gorm:"type:varchar(255)" json:"contact_info"`
	PurchaseHistory []PurchaseHistory `gorm:"foreignKey:SupplierID" json:"-"`
}

func (s Supplier) ExcelHeaders() []string {
	return []string{"ID", "NOMBRE", "INFORMACIÓN DE CONTACTO", "FECHA DE CREACIÓN"}
}

func (s Supplier) ExcelRow() []interface{} {
	return []interface{}{s.ID, s.Name, s.ContactInfo, s.CreatedAt.Format("2006-01-02 15:04:05")}
}
