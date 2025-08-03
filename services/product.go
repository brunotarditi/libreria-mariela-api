package services

import (
	"fmt"
	"io"
	"libreria/common"
	"libreria/models"
	"libreria/repositories"
	"libreria/responses"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProductsWithCategoriesAndBrands() ([]responses.ProductResponse, error)
	ExportToExcel() (*excelize.File, error)
	ImportFromExcel(reader io.Reader) error
}

type productService struct {
	db          *gorm.DB
	productRepo repositories.ProductRepository
	categoryOps *common.GormOperations[models.Category]
	brandOps    *common.GormOperations[models.Brand]
}

func NewProductService(db *gorm.DB, productRepo repositories.ProductRepository, categoryOps *common.GormOperations[models.Category], brandOps *common.GormOperations[models.Brand]) ProductService {
	return &productService{
		db:          db,
		productRepo: productRepo,
		categoryOps: categoryOps,
		brandOps:    brandOps,
	}
}

func (s *productService) GetAllProductsWithCategoriesAndBrands() ([]responses.ProductResponse, error) {
	products, err := s.productRepo.FindAll()

	if err != nil {
		return []responses.ProductResponse{}, err
	}

	return products, nil
}

func (s *productService) ExportToExcel() (*excelize.File, error) {

	f := excelize.NewFile()
	sheet := "Products"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"CÓDIGO", "SKU", "NOMBRE", "MARGEN DE GANANCIA (%)", "DESCRIPCION", "CATEGORÍA", "MARCA"}

	f.SetColWidth(sheet, "A", "B", 10)
	f.SetColWidth(sheet, "C", "E", 30)
	f.SetColWidth(sheet, "F", "G", 15)

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#576CBC"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear estilo para encabezados: %v", err)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, style)

		if h == "SKU" {
			f.AddComment(sheet, excelize.Comment{
				Author: "Sistema",
				Text:   "Código alfanumérico único asignado a cada producto o variante de producto dentro de un inventario",
				Cell:   cell,
			})
		}

		if h == "CÓDIGO" || h == "NOMBRE" || h == "MARGEN DE GANANCIA (%)" || h == "CATEGORÍA" || h == "MARCA" {
			f.AddComment(sheet, excelize.Comment{
				Author: "Sistema",
				Text:   "Campo obligatorio",
				Cell:   cell,
			})
		}
	}

	productExample := []interface{}{
		"ABC123",
		"SKU001",
		"Lapicera azul trazo fino",
		25,
		"Lapicera trazo fino de color azul",
		"Lapiceras",
		"Bic",
	}
	f.SetSheetRow(sheet, "A2", &productExample)

	_, _ = f.NewSheet("Categories")

	categories, err := s.categoryOps.Pluck("name")

	if err != nil {
		return nil, fmt.Errorf("error al obtener categorías: %v", err)
	}

	for i, category := range categories {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue("Categories", cell, category)
	}

	_, _ = f.NewSheet("Brands")

	brands, err := s.brandOps.Pluck("name")

	if err != nil {
		return nil, fmt.Errorf("error al obtener marcas: %v", err)
	}

	for i, brand := range brands {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue("Brands", cell, brand)
	}

	for row := 2; row <= 101; row++ {
		categoryCell, _ := excelize.CoordinatesToCellName(6, row) // Columna F (Category)
		brandCell, _ := excelize.CoordinatesToCellName(7, row)    // Columna G (Brand)

		dv1 := excelize.NewDataValidation(true)
		dv1.SetSqref(categoryCell)
		dv1.SetSqrefDropList("Categories!$A$1:$A$100")
		if err = f.AddDataValidation(sheet, dv1); err != nil {
			return nil, fmt.Errorf("error en validación categoría fila %d: %v", row, err)
		}

		dv2 := excelize.NewDataValidation(true)
		dv2.SetSqref(brandCell)
		dv2.SetSqrefDropList("Brands!$A$1:$A$100")

		if err = f.AddDataValidation(sheet, dv2); err != nil {
			return nil, fmt.Errorf("error en validación marca fila %d: %v", row, err)
		}
	}

	// Ocultar hojas auxiliares
	_ = f.SetSheetVisible("Categories", false)
	_ = f.SetSheetVisible("Brands", false)

	// Set "Products" como hoja activa
	index, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(index)

	return f, nil
}

func (s *productService) ImportFromExcel(reader io.Reader) error {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %v", err)
	}

	rows, err := f.GetRows("Products")
	if err != nil {
		return fmt.Errorf("no se pudo leer la hoja 'Products': %v", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("el archivo no contiene datos")
	}

	categoryMap := map[string]uint{}
	categories, _ := s.categoryOps.FindAll()
	for _, c := range categories {
		categoryMap[c.Name] = c.ID
	}

	brandMap := map[string]uint{}
	brands, _ := s.brandOps.FindAll()
	for _, b := range brands {
		brandMap[b.Name] = b.ID
	}

	var validProducts []models.Product
	var errors []string

	for i, row := range rows[1:] { // saltar encabezado (fila 0)
		rowNum := i + 2 // fila real en Excel

		if len(row) < 7 {
			errors = append(errors, fmt.Sprintf("Fila %d: Faltan columnas", rowNum))
			continue
		}

		code := strings.TrimSpace(row[0])
		sku := strings.TrimSpace(row[1])
		name := strings.TrimSpace(row[2])
		profitMarginStr := strings.TrimSpace(row[3])
		description := strings.TrimSpace(row[4])
		categoryName := strings.TrimSpace(row[5])
		brandName := strings.TrimSpace(row[6])

		// Validaciones básicas
		if code == "" || name == "" || profitMarginStr == "" || categoryName == "" || brandName == "" {
			errors = append(errors, fmt.Sprintf("Fila %d: Campos obligatorios faltantes", rowNum))
			continue
		}

		profitMargin, err := strconv.ParseFloat(profitMarginStr, 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Fila %d: Margen de ganancia inválido", rowNum))
			continue
		}

		if profitMargin < 0 || profitMargin > 100 {
			errors = append(errors, fmt.Sprintf("Fila %d: Margen de ganancia debe estar entre 0 y 100", rowNum))
			continue
		}

		if exists, err := s.productRepo.ExistsByCodeAndName(code, name); err != nil || exists {
			errors = append(errors, fmt.Sprintf("Fila %d: Producto con código '%s' y nombre '%s' ya existe", rowNum, code, name))
			continue
		}
		if sku != "" {
			if exists, err := s.productRepo.ExistsBySku(sku); err != nil || exists {
				errors = append(errors, fmt.Sprintf("Fila %d: SKU '%s' ya existe", rowNum, sku))
				continue
			}
		}

		categoryID, ok := categoryMap[categoryName]
		if !ok {
			err := fmt.Sprintf("fila %d: Categoría '%s' no encontrada", rowNum, categoryName)
			errors = append(errors, err)
			continue
		}

		brandID, ok := brandMap[brandName]
		if !ok {
			err := fmt.Sprintf("Fila %d: Marca '%s' no encontrada", rowNum, brandName)
			errors = append(errors, err)
			continue
		}

		p := models.Product{
			Code:         code,
			Sku:          sku,
			Name:         name,
			ProfitMargin: profitMargin,
			Description:  description,
			CategoryID:   categoryID,
			BrandID:      brandID,
		}
		validProducts = append(validProducts, p)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores encontrados:\n%s", strings.Join(errors, "\n"))
	}

	if _, err := s.productRepo.CreateMany(validProducts); err != nil {
		return fmt.Errorf("error al guardar productos: %v", err)
	}

	return nil

}
