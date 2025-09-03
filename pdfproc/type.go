package pdfproc

import (
	"pdfimporter/assets"
	"pdfimporter/domain"

	"github.com/mechiko/maroto/v2/pkg/core"
)

type pdfProc struct {
	domain.Apper
	maroto             core.Maroto
	assets             *assets.Assets
	templateDatamatrix *MarkTemplate
	templateBar        *MarkTemplate
	document           core.Document
	debug              bool
}

func New(app domain.Apper, tmplDatamatrix, tmplBar *MarkTemplate, assets *assets.Assets) (*pdfProc, error) {
	p := &pdfProc{
		Apper:              app,
		templateDatamatrix: tmplDatamatrix,
		templateBar:        tmplBar,
		assets:             assets,
	}
	// if err := p.BuildMaroto(); err != nil {
	// 	return nil, fmt.Errorf("build maroto error %w", err)
	// }
	return p, nil
}
