package pdfproc

import (
	"fmt"
	"pdfimporter/embeded"

	"github.com/mechiko/maroto/v2"

	"github.com/mechiko/maroto/v2/pkg/config"
	"github.com/mechiko/maroto/v2/pkg/consts/fontstyle"
	"github.com/mechiko/maroto/v2/pkg/repository"

	"github.com/mechiko/maroto/v2/pkg/props"
)

func (p *pdfProc) PdfDocumentSave(fileName string) (err error) {
	if p.document == nil {
		return fmt.Errorf("save document is nil %w", err)
	}
	err = p.document.Save(fileName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (p *pdfProc) PdfDocumentReportSave(fileName string) (err error) {
	// запись отчета генерации
	if p.document == nil {
		return fmt.Errorf("save document report is nil %w", err)
	}
	err = p.document.GetReport().Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}
	return nil
}

func (p *pdfProc) BuildMaroto(width, height float64) (err error) {
	customFont := "arial-unicode-ms"
	customFonts, err := repository.New().
		AddUTF8FontFromBytes(customFont, fontstyle.Normal, embeded.Regular).
		AddUTF8FontFromBytes(customFont, fontstyle.Italic, embeded.Italic).
		AddUTF8FontFromBytes(customFont, fontstyle.Bold, embeded.Bold).
		AddUTF8FontFromBytes(customFont, fontstyle.BoldItalic, embeded.BoldItalic).
		Load()
	if err != nil {
		return err
	}
	// …custom font setup…
	builder := config.NewBuilder().WithCustomFonts(customFonts)
	cfg := builder.WithDefaultFont(&props.Font{Family: customFont}).Build()
	cfg.Dimensions.Height = height
	cfg.Dimensions.Width = width
	cfg.Margins.Bottom = 0
	cfg.Margins.Top = 1
	cfg.Margins.Left = 2
	cfg.Margins.Right = 2
	cfg.DefaultFont.Size = 4
	p.maroto = maroto.New(cfg)
	p.maroto = maroto.NewMetricsDecorator(p.maroto)
	return nil
}

func (p *pdfProc) DocumentGenerate() (err error) {
	doc, err := p.maroto.Generate()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	p.document = doc
	return nil
}

func (p *pdfProc) AddPageByTemplate(tmpl *MarkTemplate, kod string, ser string) error {
	if tmpl == nil {
		return fmt.Errorf("add page template is nil")
	}
	pgNew, err := p.Page(tmpl, kod, ser)
	if err != nil {
		return fmt.Errorf("page error %w", err)
	}
	if pgNew == nil {
		return fmt.Errorf("page nil")
	}
	p.maroto.AddPages(pgNew)
	return nil
}
