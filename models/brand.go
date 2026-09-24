package models

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(65);not null" json:"name"`
	Products []Product `gorm:"foreignKey:BrandID" json:"-"`
}

func (b Brand) ExcelHeaders() []string {
	return []string{"ID", "NOMBRE", "FECHA DE CREACIÓN"}
}

func (b Brand) ExcelRow() []interface{} {
	return []interface{}{b.ID, b.Name, b.CreatedAt.Format("2006-01-02 15:04:05")}
}
