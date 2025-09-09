package pdfproc

import (
	"github.com/mechiko/maroto/v2/pkg/consts/border"
	"github.com/mechiko/maroto/v2/pkg/consts/linestyle"
	"github.com/mechiko/maroto/v2/pkg/props"
)

var colStyle = &props.Cell{
	BackgroundColor: &props.Color{Red: 128, Green: 128, Blue: 128},
	BorderType:      border.None,
	BorderColor:     &props.Color{Red: 200, Green: 0, Blue: 0},
	LineStyle:       linestyle.Dashed,
	BorderThickness: 0.5,
}
