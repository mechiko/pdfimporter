package pdfproc

import (
	"fmt"
	"pdfimporter/assets"
	"pdfimporter/domain"

	"github.com/mechiko/maroto/v2/pkg/core"
)

type pdfProc struct {
	domain.Apper
	maroto   core.Maroto
	assets   *assets.Assets
	document core.Document
	debug    bool
	height   float64
	width    float64
}

func New(app domain.Apper, assets *assets.Assets) (*pdfProc, error) {
	if app == nil {
		return nil, fmt.Errorf("app is nil")
	}
	if assets == nil {
		return nil, fmt.Errorf("assets is nil")
	}
	p := &pdfProc{
		Apper:  app,
		assets: assets,
	}
	return p, nil
}
