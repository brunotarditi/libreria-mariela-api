package controllers

import (
	"encoding/json"
	"libreria/responses"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockSearchService struct{}

func (m *mockSearchService) Search(query string, limit ...int) (responses.GlobalSearchResponse, error) {
	return responses.GlobalSearchResponse{
		Products: []responses.SearchItem{
			{ID: 1, Title: "Cuaderno Rivadavia", Subtitle: "Código: C-102", Route: "/products/detail/1"},
		},
		Brands: []responses.SearchItem{
			{ID: 4, Title: "Bic", Subtitle: "Marca de librería", Route: "/brands/detail/4"},
		},
		Categories: []responses.SearchItem{
			{ID: 2, Title: "Escolares", Subtitle: "Categoría", Route: "/categories/detail/2"},
		},
		Suppliers: []responses.SearchItem{
			{ID: 3, Title: "Lito Distribuidora", Subtitle: "Contacto: 112345678", Route: "/suppliers/detail/3"},
		},
		Customers: []responses.SearchItem{
			{ID: 8, Title: "Juan Pérez", Subtitle: "Cliente frecuente", Route: "/customers/detail/8"},
		},
	}, nil
}

func TestSearchController_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewSearchController(&mockSearchService{})

	r := gin.New()
	r.GET("/search", ctrl.Search)

	req, _ := http.NewRequest(http.MethodGet, "/search?q=cuaderno", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var res responses.GlobalSearchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(res.Products) != 1 || res.Products[0].Title != "Cuaderno Rivadavia" {
		t.Errorf("unexpected products in search response: %+v", res.Products)
	}
	if len(res.Brands) != 1 || res.Brands[0].Title != "Bic" {
		t.Errorf("unexpected brands in search response: %+v", res.Brands)
	}
}
