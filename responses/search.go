package responses

type SearchItem struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Route    string `json:"route"`
}

type GlobalSearchResponse struct {
	Products   []SearchItem `json:"products"`
	Brands     []SearchItem `json:"brands"`
	Categories []SearchItem `json:"categories"`
	Suppliers  []SearchItem `json:"suppliers"`
	Customers  []SearchItem `json:"customers"`
}
