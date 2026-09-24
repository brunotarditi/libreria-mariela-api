package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type mockOperations struct {
	deleteManyFunc func(ids []uint) (int64, error)
	findAllFunc    func() ([]DummyModel, error)
}

func (m *mockOperations) FindAll() ([]DummyModel, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return nil, nil
}
func (m *mockOperations) FindByID(id string) (DummyModel, error)                { return DummyModel{}, nil }
func (m *mockOperations) Paginated(offset, size int, options QueryOptions) ([]DummyModel, error) {
	return nil, nil
}
func (m *mockOperations) Count(options QueryOptions) (int64, error)             { return 0, nil }
func (m *mockOperations) Create(model DummyModel) (DummyModel, error)           { return model, nil }
func (m *mockOperations) CreateMany(model []DummyModel) error                   { return nil }
func (m *mockOperations) Update(model DummyModel) (DummyModel, error)           { return model, nil }
func (m *mockOperations) Delete(id string) error                                { return nil }
func (m *mockOperations) DeleteMany(ids []uint) (int64, error) {
	if m.deleteManyFunc != nil {
		return m.deleteManyFunc(ids)
	}
	return int64(len(ids)), nil
}
func (m *mockOperations) Pluck(field string) ([]string, error)                  { return nil, nil }

func TestBulkDeleteHandler_Success_Multiple(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOps := &mockOperations{
		deleteManyFunc: func(ids []uint) (int64, error) {
			return int64(len(ids)), nil
		},
	}

	r := gin.New()
	r.POST("/bulk-delete", BulkDelete[DummyModel](mockOps))

	reqBody := []byte(`{"ids": [1, 2, 5, 10]}`)
	req, _ := http.NewRequest(http.MethodPost, "/bulk-delete", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success: true, got %v", res["success"])
	}
	if res["deleted_count"] != float64(4) {
		t.Errorf("expected deleted_count: 4, got %v", res["deleted_count"])
	}
	if res["message"] != "Se eliminaron 4 registros con éxito" {
		t.Errorf("expected message 'Se eliminaron 4 registros con éxito', got %q", res["message"])
	}
}

func TestBulkDeleteHandler_Success_Single(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOps := &mockOperations{
		deleteManyFunc: func(ids []uint) (int64, error) {
			return 1, nil
		},
	}

	r := gin.New()
	r.POST("/bulk-delete", BulkDelete[DummyModel](mockOps))

	reqBody := []byte(`{"ids": [1]}`)
	req, _ := http.NewRequest(http.MethodPost, "/bulk-delete", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success: true, got %v", res["success"])
	}
	if res["deleted_count"] != float64(1) {
		t.Errorf("expected deleted_count: 1, got %v", res["deleted_count"])
	}
	if res["message"] != "Se eliminó 1 registro con éxito" {
		t.Errorf("expected message 'Se eliminó 1 registro con éxito', got %q", res["message"])
	}
}

func TestBulkDeleteHandler_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOps := &mockOperations{}

	r := gin.New()
	r.POST("/bulk-delete", BulkDelete[DummyModel](mockOps))

	tests := []struct {
		name    string
		body    string
		expCode int
	}{
		{
			name:    "empty body",
			body:    ``,
			expCode: http.StatusBadRequest,
		},
		{
			name:    "empty ids array",
			body:    `{"ids": []}`,
			expCode: http.StatusBadRequest,
		},
		{
			name:    "contains zero ID",
			body:    `{"ids": [0]}`,
			expCode: http.StatusBadRequest,
		},
		{
			name:    "non-numeric ID in array",
			body:    `{"ids": ["abc"]}`,
			expCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/bulk-delete", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expCode {
				t.Errorf("expected code %d, got %d, body: %s", tt.expCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestBulkDeleteHandler_ServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOps := &mockOperations{
		deleteManyFunc: func(ids []uint) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	r := gin.New()
	r.POST("/bulk-delete", BulkDelete[DummyModel](mockOps))

	reqBody := []byte(`{"ids": [1, 2]}`)
	req, _ := http.NewRequest(http.MethodPost, "/bulk-delete", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestExportHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOps := &mockOperations{
		findAllFunc: func() ([]DummyModel, error) {
			return []DummyModel{
				{ID: 1, Name: "Item 1"},
				{ID: 2, Name: "Item 2"},
			}, nil
		},
	}

	r := gin.New()
	r.GET("/export", Export[DummyModel](mockOps, "dummies"))

	req, _ := http.NewRequest(http.MethodGet, "/export", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	expectedContentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if w.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, w.Header().Get("Content-Type"))
	}

	expectedDisposition := `attachment; filename="dummies_export.xlsx"`
	if w.Header().Get("Content-Disposition") != expectedDisposition {
		t.Errorf("expected Content-Disposition %q, got %q", expectedDisposition, w.Header().Get("Content-Disposition"))
	}

	reader, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("failed to open returned Excel buffer: %v", err)
	}
	defer reader.Close()

	header1, _ := reader.GetCellValue("Dummies", "A1")
	header2, _ := reader.GetCellValue("Dummies", "B1")
	if header1 != "ID" || header2 != "NOMBRE" {
		t.Errorf("expected headers ID, NOMBRE, got %q, %q", header1, header2)
	}

	val1, _ := reader.GetCellValue("Dummies", "A2")
	val2, _ := reader.GetCellValue("Dummies", "B2")
	if val1 != "1" || val2 != "Item 1" {
		t.Errorf("expected row 2 to be 1, Item 1, got %q, %q", val1, val2)
	}
}
