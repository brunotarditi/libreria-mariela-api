package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func TestBuildExcel(t *testing.T) {
	headers := []string{"ID", "NOMBRE", "PRECIO"}
	rows := [][]interface{}{
		{uint(1), "Cuaderno", 1500.50},
		{uint(2), "Lapicera", 250.00},
	}

	f, err := BuildExcel("Productos", headers, rows)
	if err != nil {
		t.Fatalf("BuildExcel failed: %v", err)
	}

	// Verify sheet exists
	index, err := f.GetSheetIndex("Productos")
	if err != nil || index == -1 {
		t.Fatalf("expected sheet 'Productos' to exist, got index %d, err %v", index, err)
	}

	// Verify header cell
	val, err := f.GetCellValue("Productos", "A1")
	if err != nil || val != "ID" {
		t.Errorf("expected A1 to be 'ID', got %q, err %v", val, err)
	}

	// Verify data cell
	val, err = f.GetCellValue("Productos", "B2")
	if err != nil || val != "Cuaderno" {
		t.Errorf("expected B2 to be 'Cuaderno', got %q, err %v", val, err)
	}

	// Verify column width is set (greater than 0)
	colWidth, err := f.GetColWidth("Productos", "B")
	if err != nil || colWidth <= 0 {
		t.Errorf("expected column B width > 0, got %f, err %v", colWidth, err)
	}
}

func TestStreamExcel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	f, err := BuildExcel("TestSheet", []string{"Col1"}, [][]interface{}{{"Val1"}})
	if err != nil {
		t.Fatalf("BuildExcel failed: %v", err)
	}

	r := gin.New()
	r.GET("/export", func(c *gin.Context) {
		if err := StreamExcel(c, f, "test_export.xlsx"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	req, _ := http.NewRequest(http.MethodGet, "/export", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	expectedContentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if w.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, w.Header().Get("Content-Type"))
	}

	expectedDisposition := `attachment; filename="test_export.xlsx"`
	if w.Header().Get("Content-Disposition") != expectedDisposition {
		t.Errorf("expected Content-Disposition %q, got %q", expectedDisposition, w.Header().Get("Content-Disposition"))
	}

	// Verify returned body is a valid excel file
	reader, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("expected response body to be valid Excel file, failed to open: %v", err)
	}
	defer reader.Close()

	val, err := reader.GetCellValue("TestSheet", "A1")
	if err != nil || val != "Col1" {
		t.Errorf("expected A1 in parsed file to be 'Col1', got %q, err %v", val, err)
	}
}
