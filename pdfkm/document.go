package pdfkm

import (
	"fmt"
	"path/filepath"
	"pdfimporter/domain/models/application"
	"pdfimporter/pdfproc"
	"slices"
	"strings"
	"time"

	"github.com/mechiko/utility"
)

func (k *Pdf) Document(model *application.Application, ch chan float64) (string, error) {
	pdfDocument, err := pdfproc.New(k, k.templateDatamatrix, k.templateBar, k.assets)
	if err != nil {
		return "", fmt.Errorf("Error create pdfproc: %v", err)
	}
	if err := pdfDocument.BuildMaroto(float64(model.Width), float64(model.Height)); err != nil {
		return "", fmt.Errorf("build maroto error %w", err)
	}

	start := time.Now()

	totalItems := len(k.Cis) + len(k.Pallet)
	step := 0.0
	if totalItems > 0 {
		step = 99.0 / float64(totalItems)
	}
	palets := make([]string, 0, len(k.Pallet))
	for k2 := range k.Pallet {
		palets = append(palets, k2)
	}
	slices.Sort(palets)
	i := 0
	for _, palet := range palets {
		cises := k.Pallet[palet]
		for _, cis := range cises {
			fnc := cis.FNC1()
			ser := cis.Serial
			pdfDocument.AddPageByTemplate(k.templateDatamatrix, fnc, ser)
			ch <- step * float64(i)
			i++
		}
		pdfDocument.AddPageByTemplate(k.templateBar, palet, "")
		ch <- step * float64(i)
		i++
		if k.DebugMode() {
			if i > 26 {
				break
			}
		}
	}
	k.Logger().Debugf("заняло времени %s", time.Since(start))

	k.Logger().Debugf("buid pages %v", time.Since(start))
	start = time.Now()
	err = pdfDocument.DocumentGenerate()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	k.Logger().Debugf("generate document %v\n", time.Since(start))

	fileName := "PDF_" + filepath.Base(model.File)
	fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
	fileName = utility.TimeFileName(fileName) + ".pdf"
	filePdf, err := utility.DialogSaveFile(utility.Pdf, fileName, ".")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if filePdf != "" {
		if !strings.HasSuffix(filePdf, ".pdf") {
			filePdf = fmt.Sprintf("%s.pdf", filePdf)
		}
		err = pdfDocument.PdfDocumentSave(filePdf)
		if err != nil {
			return "", err
		}
		ch <- 100.00
	}
	// запись отчета генерации
	err = pdfDocument.PdfDocumentReportSave("report.txt")
	if err != nil {
		return "", fmt.Errorf("failed to save report: %w", err)
	}
	return filePdf, nil
}
