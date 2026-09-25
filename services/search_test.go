package services

import (
	"libreria/models"
	"testing"
)

type mockSearchRepo struct{}

func (m *mockSearchRepo) SearchProducts(query string, limit int) ([]models.Product, error) {
	return []models.Product{
		{
			Code: "C-102",
			Sku:  "7791234",
			Name: "Cuaderno Rivadavia",
		},
	}, nil
}

func (m *mockSearchRepo) SearchBrands(query string, limit int) ([]models.Brand, error) {
	return []models.Brand{
		{
			Name: "Bic",
		},
	}, nil
}

func (m *mockSearchRepo) SearchCategories(query string, limit int) ([]models.Category, error) {
	return []models.Category{
		{
			Name: "Escolares",
		},
	}, nil
}

func (m *mockSearchRepo) SearchSuppliers(query string, limit int) ([]models.Supplier, error) {
	return []models.Supplier{
		{
			Name:        "Lito Distribuidora",
			ContactInfo: "112345678",
		},
	}, nil
}

func (m *mockSearchRepo) SearchCustomers(query string, limit int) ([]models.Customer, error) {
	return []models.Customer{
		{
			Name:        "Juan Pérez",
			ContactInfo: "119876543",
		},
	}, nil
}

func TestSearchService_EmptyQuery(t *testing.T) {
	service := NewSearchService(nil, &mockSearchRepo{})

	res, err := service.Search("   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Products) != 0 || len(res.Brands) != 0 || len(res.Categories) != 0 || len(res.Suppliers) != 0 || len(res.Customers) != 0 {
		t.Errorf("expected all empty slices for blank query, got %+v", res)
	}
}

func TestSearchService_ValidQuery(t *testing.T) {
	service := NewSearchService(nil, &mockSearchRepo{})

	res, err := service.Search("cuaderno")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(res.Products))
	}
	if res.Products[0].Title != "Cuaderno Rivadavia" {
		t.Errorf("expected product title 'Cuaderno Rivadavia', got %s", res.Products[0].Title)
	}
	if res.Brands[0].Subtitle != "Marca de librería" {
		t.Errorf("expected brand subtitle 'Marca de librería', got %s", res.Brands[0].Subtitle)
	}
	if res.Categories[0].Subtitle != "Categoría" {
		t.Errorf("expected category subtitle 'Categoría', got %s", res.Categories[0].Subtitle)
	}
	if res.Suppliers[0].Subtitle != "Contacto: 112345678" {
		t.Errorf("expected supplier subtitle 'Contacto: 112345678', got %s", res.Suppliers[0].Subtitle)
	}
	if res.Customers[0].Subtitle != "Contacto: 119876543" {
		t.Errorf("expected customer subtitle 'Contacto: 119876543', got %s", res.Customers[0].Subtitle)
	}
}
