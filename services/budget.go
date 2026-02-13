package services

import (
	"log"

	"github.com/johnfercher/maroto/v2"
	"gorm.io/gorm"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"

	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type BudgetService interface {
	GeneratePDF() ([]byte, error)
}

type budgetService struct {
	db *gorm.DB
}

func NewBudgetService(db *gorm.DB) BudgetService {
	return &budgetService{
		db: db,
	}
}

func (s *budgetService) GeneratePDF() ([]byte, error) {

	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	darkGrayColor := getDarkGrayColor()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())
	if err != nil {
		log.Fatal(err.Error())
	}

	err = m.RegisterFooter(getPageFooter())
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(text.NewRow(10, "Presupuesto", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  14,
	}),
		text.NewRow(10, "Cliente: Fito", props.Text{Align: align.Left}),
		text.NewRow(10, "Descripción: Útiles para primero de Pompeya", props.Text{Align: align.Left}),
		text.NewRow(10, "Fecha: 07/08/2025", props.Text{Align: align.Left}),
		text.NewRow(10, "Vence: 14/08/2025", props.Text{Align: align.Left}),
	)

	m.AddRow(7,
		text.NewCol(3, "Lista de productos", props.Text{
			Top:   1.5,
			Size:  9,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})

	m.AddRows(getItems()...)

	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
	return document.GetBytes(), nil
}

func getItems() []core.Row {
	rows := []core.Row{
		row.New(4).Add(
			text.NewCol(5, "Producto", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Cantidad", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(3, "Precio unitario", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Subtotal", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
		),
	}

	var contentsRow []core.Row
	contents := getContents()

	for i, content := range contents {
		r := row.New(4).Add(
			text.NewCol(5, content[0], props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, content[1], props.Text{Size: 8, Align: align.Center}),
			text.NewCol(3, content[2], props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, content[3], props.Text{Size: 8, Align: align.Center}),
		)
		if i%2 == 0 {
			gray := getGrayColor()
			r.WithStyle(&props.Cell{BackgroundColor: gray})
		}

		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	rows = append(rows, row.New(20).Add(
		col.New(6),
		text.NewCol(2, "Total:", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Right,
		}),
		text.NewCol(3, "$ 22.800,00", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Center,
		}),
	))

	return rows
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(3, "assets/public/logo.png", props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(6),
		col.New(3).Add(
			text.New("Paso de los Andes 902, M5501 Godoy Cruz, Mendoza, Argentina.", props.Text{
				Size:  8,
				Align: align.Right,
				Color: getPrimaryColor(),
			}),
			text.New("Tel: +54 9 2613 37-8438", props.Text{
				Top:   12,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: getBlueColor(),
			}),
			text.New("https://mariela.pages.dev", props.Text{
				Top:   15,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: getBlueColor(),
			}),
		),
	)
}

func getPageFooter() core.Row {
	return row.New(20).Add(
		col.New(12).Add(
			text.New("Tel: +54 9 2613 37-8438", props.Text{
				Top:   13,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Left,
				Color: getBlueColor(),
			}),
			text.New("https://mariela.pages.dev/", props.Text{
				Top:   16,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Left,
				Color: getBlueColor(),
			}),
		),
	)
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   6,
		Green: 6,
		Blue:  173,
	}
}

func getPrimaryColor() *props.Color {
	return &props.Color{
		Red:   35,
		Green: 155,
		Blue:  167,
	}
}

func getContents() [][]string {
	return [][]string{
		{"Lapicera Bic", "2", "$ 800,00", "$ 1.600,00"},
		{"Hojas Rivadavia A4", "3", "$ 1.700,00", "$ 5.100,00"},
		{"Tijera Mapped", "1", "$ 1200,00", "$ 1.200,00"},
		{"Goma Mapped", "4", "$ 350,00", "$ 1.400,00"},
		{"Cuaderno Rivadavia", "3", "$ 4.500,00", "$ 13.500,00"},
	}
}
