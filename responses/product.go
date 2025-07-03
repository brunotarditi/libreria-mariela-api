package responses

type ProductResponse struct {
	ID           uint    `json:"id"`
	Code         string  `json:"code"`
	Sku          string  `json:"sku"`
	Name         string  `json:"name"`
	ProfitMargin float64 `json:"profit_margin"`
	Description  string  `json:"description"`
	CategoryName string  `json:"category_name"`
	BrandName    string  `json:"brand_name"`
}
