package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Code              string            `gorm:"type:varchar(20);not null" json:"code"`
	Sku               string            `gorm:"type:varchar(20);not null" json:"sku"`
	Name              string            `gorm:"type:varchar(65);not null" json:"name"`
	ProfitMargin      float64           `gorm:"type:decimal(5,2);not null" json:"profit_margin"`
	Description       string            `gorm:"type:varchar(150);not null" json:"description"`
	CategoryID        uint              `gorm:"not null" json:"category_id"`
	Category          Category          `gorm:"foreignKey:CategoryID" json:"category"`
	BrandID           uint              `gorm:"not null" json:"brand_id"`
	Brand             Brand             `gorm:"foreignKey:BrandID" json:"brand"`
	StockMovements    []StockMovement   `gorm:"foreignKey:ProductID" json:"-"`
	ProductStocks     []ProductStock    `gorm:"foreignKey:ProductID" json:"-"`
	PriceLists        []PriceList       `gorm:"foreignKey:ProductID" json:"-"`
	PurchaseHistories []PurchaseHistory `gorm:"foreignKey:ProductID" json:"-"`
	SellHistories     []SellHistory     `gorm:"foreignKey:ProductID" json:"-"`
}
