package pdfproc

import (
	"github.com/mechiko/maroto/v2/pkg/consts/align"
	"github.com/mechiko/maroto/v2/pkg/consts/barcode"
	"github.com/mechiko/maroto/v2/pkg/consts/border"
	"github.com/mechiko/maroto/v2/pkg/consts/extension"
	"github.com/mechiko/maroto/v2/pkg/consts/fontstyle"
	"github.com/mechiko/maroto/v2/pkg/consts/linestyle"
	"github.com/mechiko/maroto/v2/pkg/props"
)

// описывает часть строки
type RowPrimitive struct {
	RowHeight  float64        `json:"row_height,omitempty"`
	ColWidth   int            `json:"col_width,omitempty"`
	Value      string         `json:"value,omitempty"`
	Style      fontstyle.Type `json:"style,omitempty"`
	FontSize   float64        `json:"font_size,omitempty"`
	Align      align.Type     `json:"align,omitempty"`
	Bar        string         `json:"bar,omitempty"`
	DataMatrix string         `json:"data_matrix,omitempty"`
	// Left is the space between the left cell boundary to the rectangle, if center is false.
	Left float64 `json:"left,omitempty"`
	// Top is space between the upper cell limit to the barcode, if center is false.
	Top float64 `json:"top,omitempty"`
	// Percent is how much the rectangle will occupy the cell,
	// ex 100%: The rectangle will fulfill the entire cell
	// ex 50%: The greater side from the rectangle will have half the size of the cell.
	Percent float64 `json:"percent,omitempty"`
	// indicate whether only the width should be used as a reference to calculate the component size, disregarding the height
	// ex true: The component will be scaled only based on the available width, disregarding the available height
	JustReferenceWidth bool `json:"just_reference_width,omitempty"`
	// Center define that the barcode will be vertically and horizontally centralized.
	Center bool `json:"center,omitempty"`
	// Proportion is the proportion between size of the barcode.
	// Ex: 16x9, 4x3...
	Proportion props.Proportion `json:"proportion,omitempty"`
	// Center define that the barcode will be vertically and horizontally centralized.
	Type       barcode.Type   `json:"type,omitempty"`
	Image      string         `json:"image,omitempty"`
	ImageExt   extension.Type `json:"image_ext,omitempty"`
	ImageDebug bool           `json:"image_debug,omitempty"`
}

func (p *RowPrimitive) PropsText() props.Text {
	return props.Text{
		Top:   p.Top,
		Style: p.Style,
		Size:  p.FontSize,
		Align: p.Align,
	}
}

func (p *RowPrimitive) PropsRect() props.Rect {
	return props.Rect{
		Left:               p.Left,
		Top:                p.Top,
		Percent:            p.Percent,
		Center:             p.Center,
		JustReferenceWidth: p.JustReferenceWidth,
	}
}

func (p *RowPrimitive) PropsBar() props.Barcode {
	return props.Barcode{
		Left:       p.Left,
		Top:        p.Top,
		Percent:    p.Percent,
		Center:     p.Center,
		Proportion: p.Proportion,
		Type:       p.Type,
	}
}

var colStyle = &props.Cell{
	BackgroundColor: &props.Color{Red: 128, Green: 128, Blue: 128},
	BorderType:      border.None,
	BorderColor:     &props.Color{Red: 200, Green: 0, Blue: 0},
	LineStyle:       linestyle.Dashed,
	BorderThickness: 0.5,
}

// шаблон состоит из строк упорядоченных по ключу map
// если примитив один значит целиком одна строка, может быть с колонкой
// если примитивов больше одного то это строка из колонок
// если RowHeight 0 то это text.NewAutoRow
// если Value пусто это строка выравнивания с высотой
type MarkTemplate struct {
	Name       string
	PageWidth  float64
	PageHeight float64
	Rows       map[string][]*RowPrimitive
}
