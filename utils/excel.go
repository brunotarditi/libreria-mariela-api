package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

const (
	CorporateColor = "#239BA7"
)

// BuildExcel creates a styled excelize.File with corporate headers and auto-adjusted column widths.
func BuildExcel(sheetName string, headers []string, rows [][]interface{}) (*excelize.File, error) {
	f := excelize.NewFile()
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	f.SetSheetName("Sheet1", sheetName)

	// Corporate header style: bold, white text, #239BA7 fill, centered, borders
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{CorporateColor},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#B0C4DE", Style: 1},
			{Type: "bottom", Color: "#B0C4DE", Style: 1},
			{Type: "left", Color: "#B0C4DE", Style: 1},
			{Type: "right", Color: "#B0C4DE", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear estilo de encabezado: %w", err)
	}

	// Data cell styles
	defaultStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#E5E5E5", Style: 1},
			{Type: "bottom", Color: "#E5E5E5", Style: 1},
			{Type: "left", Color: "#E5E5E5", Style: 1},
			{Type: "right", Color: "#E5E5E5", Style: 1},
		},
	})

	numberStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		CustomNumFmt: &[]string{"#,##0.00"}[0],
		Border: []excelize.Border{
			{Type: "top", Color: "#E5E5E5", Style: 1},
			{Type: "bottom", Color: "#E5E5E5", Style: 1},
			{Type: "left", Color: "#E5E5E5", Style: 1},
			{Type: "right", Color: "#E5E5E5", Style: 1},
		},
	})

	integerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		CustomNumFmt: &[]string{"0"}[0],
		Border: []excelize.Border{
			{Type: "top", Color: "#E5E5E5", Style: 1},
			{Type: "bottom", Color: "#E5E5E5", Style: 1},
			{Type: "left", Color: "#E5E5E5", Style: 1},
			{Type: "right", Color: "#E5E5E5", Style: 1},
		},
	})

	// Set header row height
	_ = f.SetRowHeight(sheetName, 1, 26)

	// Keep track of maximum length per column for auto-adjustment
	colMaxLens := make([]int, len(headers))

	// Write headers
	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", colName)
		_ = f.SetCellValue(sheetName, cell, header)
		_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
		colMaxLens[i] = len(header)
	}

	// Write rows
	for rIdx, row := range rows {
		rowNum := rIdx + 2
		_ = f.SetRowHeight(sheetName, rowNum, 20)

		for cIdx, val := range row {
			if cIdx >= len(headers) {
				break
			}
			colName, _ := excelize.ColumnNumberToName(cIdx + 1)
			cell := fmt.Sprintf("%s%d", colName, rowNum)

			strVal := fmt.Sprintf("%v", val)
			if len(strVal) > colMaxLens[cIdx] {
				colMaxLens[cIdx] = len(strVal)
			}

			_ = f.SetCellValue(sheetName, cell, val)

			// Apply appropriate style based on type
			switch val.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				_ = f.SetCellStyle(sheetName, cell, cell, integerStyle)
			case float32, float64:
				_ = f.SetCellStyle(sheetName, cell, cell, numberStyle)
			default:
				_ = f.SetCellStyle(sheetName, cell, cell, defaultStyle)
			}
		}
	}

	// Auto-fit column widths
	for i, maxLen := range colMaxLens {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		width := float64(maxLen) + 5
		if width < 12 {
			width = 12
		}
		if width > 60 {
			width = 60
		}
		_ = f.SetColWidth(sheetName, colName, colName, width)
	}

	return f, nil
}

// StreamExcel sends the excelize.File directly as a binary response.
func StreamExcel(c *gin.Context, f *excelize.File, filename string) error {
	buf, err := f.WriteToBuffer()
	if err != nil {
		return err
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	return nil
}

// CapitalizeFirst returns the string with the first letter capitalized.
func CapitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
