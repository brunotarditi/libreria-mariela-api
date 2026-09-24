package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(65);not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"-"`
}

func (c Category) ExcelHeaders() []string {
	return []string{"ID", "NOMBRE", "FECHA DE CREACIÓN"}
}

func (c Category) ExcelRow() []interface{} {
	return []interface{}{c.ID, c.Name, c.CreatedAt.Format("2006-01-02 15:04:05")}
}
