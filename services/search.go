package services

import (
	"fmt"
	"libreria/repositories"
	"libreria/responses"
	"strings"

	"gorm.io/gorm"
)

type SearchService interface {
	Search(query string, limit ...int) (responses.GlobalSearchResponse, error)
}

type searchService struct {
	db         *gorm.DB
	searchRepo repositories.SearchRepository
}

func NewSearchService(db *gorm.DB, repo repositories.SearchRepository) SearchService {
	return &searchService{
		db:         db,
		searchRepo: repo,
	}
}

func (s *searchService) Search(query string, limit ...int) (responses.GlobalSearchResponse, error) {
	defaultLimit := 8
	if len(limit) > 0 && limit[0] > 0 {
		defaultLimit = limit[0]
	}

	cleanQuery := strings.TrimSpace(query)
	emptyResponse := responses.GlobalSearchResponse{
		Products:   []responses.SearchItem{},
		Brands:     []responses.SearchItem{},
		Categories: []responses.SearchItem{},
		Suppliers:  []responses.SearchItem{},
		Customers:  []responses.SearchItem{},
	}

	if cleanQuery == "" {
		return emptyResponse, nil
	}

	products, err := s.searchRepo.SearchProducts(cleanQuery, defaultLimit)
	if err != nil {
		return emptyResponse, err
	}

	brands, err := s.searchRepo.SearchBrands(cleanQuery, defaultLimit)
	if err != nil {
		return emptyResponse, err
	}

	categories, err := s.searchRepo.SearchCategories(cleanQuery, defaultLimit)
	if err != nil {
		return emptyResponse, err
	}

	suppliers, err := s.searchRepo.SearchSuppliers(cleanQuery, defaultLimit)
	if err != nil {
		return emptyResponse, err
	}

	customers, err := s.searchRepo.SearchCustomers(cleanQuery, defaultLimit)
	if err != nil {
		return emptyResponse, err
	}

	result := responses.GlobalSearchResponse{
		Products:   make([]responses.SearchItem, 0, len(products)),
		Brands:     make([]responses.SearchItem, 0, len(brands)),
		Categories: make([]responses.SearchItem, 0, len(categories)),
		Suppliers:  make([]responses.SearchItem, 0, len(suppliers)),
		Customers:  make([]responses.SearchItem, 0, len(customers)),
	}

	for _, p := range products {
		sub := fmt.Sprintf("Código: %s", p.Code)
		if p.Sku != "" {
			sub = fmt.Sprintf("Código: %s | SKU: %s", p.Code, p.Sku)
		}
		result.Products = append(result.Products, responses.SearchItem{
			ID:       p.ID,
			Title:    p.Name,
			Subtitle: sub,
			Route:    fmt.Sprintf("/products/detail/%d", p.ID),
		})
	}

	for _, b := range brands {
		result.Brands = append(result.Brands, responses.SearchItem{
			ID:       b.ID,
			Title:    b.Name,
			Subtitle: "Marca de librería",
			Route:    fmt.Sprintf("/brands/detail/%d", b.ID),
		})
	}

	for _, c := range categories {
		result.Categories = append(result.Categories, responses.SearchItem{
			ID:       c.ID,
			Title:    c.Name,
			Subtitle: "Categoría",
			Route:    fmt.Sprintf("/categories/detail/%d", c.ID),
		})
	}

	for _, sup := range suppliers {
		sub := "Proveedor"
		if sup.ContactInfo != "" {
			sub = fmt.Sprintf("Contacto: %s", sup.ContactInfo)
		}
		result.Suppliers = append(result.Suppliers, responses.SearchItem{
			ID:       sup.ID,
			Title:    sup.Name,
			Subtitle: sub,
			Route:    fmt.Sprintf("/suppliers/detail/%d", sup.ID),
		})
	}

	for _, cust := range customers {
		sub := "Cliente"
		if cust.ContactInfo != "" {
			sub = fmt.Sprintf("Contacto: %s", cust.ContactInfo)
		}
		result.Customers = append(result.Customers, responses.SearchItem{
			ID:       cust.ID,
			Title:    cust.Name,
			Subtitle: sub,
			Route:    fmt.Sprintf("/customers/detail/%d", cust.ID),
		})
	}

	return result, nil
}
