package pdfkm

import (
	"fmt"
	"path/filepath"
	"pdfimporter/domain/models/application"
	"pdfimporter/pdfproc"
	"strings"
	"time"

	"github.com/mechiko/utility"
)

func (k *Pdf) Document(model *application.Application, ch chan float64) (string, error) {
	if k.templateDatamatrix == nil {
		return "", fmt.Errorf("Error pdfkm datamatrix template is nil ")
	}
	if k.templatePack == nil {
		return "", fmt.Errorf("Error pdfkm pack template is nil ")
	}
	pdfDocument, err := pdfproc.New(k, k.assets)
	if err != nil {
		return "", fmt.Errorf("Error create pdfproc: %v", err)
	}
	if err := pdfDocument.BuildMaroto(k.templateDatamatrix.PageWidth, k.templateDatamatrix.PageHeight); err != nil {
		return "", fmt.Errorf("build maroto error %w", err)
	}

	start := time.Now()

	totalItems := len(k.Cis) + len(k.Pallet)
	step := 0.0
	if totalItems > 0 {
		step = 99.0 / float64(totalItems)
	}
	// packs := make([]string, 0, len(k.Pallet))
	// for k2 := range k.Pallet {
	// 	packs = append(packs, k2)
	// }
	// slices.Sort(packs)
	i := 0
	idxCis := 0
	idxKigu := 0
	for _, pack := range k.PackOrder {
		cises := k.Pallet[pack]
		for _, cis := range cises {
			// генерируем КМ
			fnc := cis.FNC1()
			ser := cis.Serial
			if err := pdfDocument.AddPageByTemplate(k.templateDatamatrix, fnc, ser, fmt.Sprintf("%06d", idxCis+1)); err != nil {
				return "", fmt.Errorf("add datamatrix page (idx %d): %w", i+1, err)
			}
			i++
			idxCis++
			if ch != nil {
				ch <- step * float64(i+1)
			}
		}
		if _, ok := k.Packs[pack]; !ok {
			return "", fmt.Errorf("missing KIGU for pack %s", pack)
		}
		kigu := k.Packs[pack]
		fnc := kigu.FNC1()
		ser := kigu.Serial
		// генерируем упаковку
		if err := pdfDocument.AddPageByTemplate(k.templatePack, fnc, ser, fmt.Sprintf("%06d", idxKigu+1)); err != nil {
			return "", fmt.Errorf("add pack page (pack %s, idx %d): %w", pack, idxKigu, err)
		}
		idxKigu++
		i++
		if ch != nil {
			ch <- step * float64(i)
		}
		// в режиме отладки берем только 26 знаков если их больше
		if k.DebugMode() {
			// в отладке генерируем две палеты
			if idxKigu > 1 {
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

	fileName := "PDF_" + filepath.Base(model.FileCIS)
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
			return "", fmt.Errorf("save pdf %q: %w", filePdf, err)
		}
		if ch != nil {
			ch <- 100.00
		}
	}
	// запись отчета генерации
	err = pdfDocument.PdfDocumentReportSave("report.txt")
	if err != nil {
		return "", fmt.Errorf("failed to save report: %w", err)
	}
	return filePdf, nil
}

func (k *Pdf) DocumentWithoutPack(model *application.Application, ch chan float64) (string, error) {
	if k.templateDatamatrix == nil {
		return "", fmt.Errorf("Error pdfkm datamatrix template is nil ")
	}
	pdfDocument, err := pdfproc.New(k, k.assets)
	if err != nil {
		return "", fmt.Errorf("Error create pdfproc: %v", err)
	}
	if err := pdfDocument.BuildMaroto(k.templateDatamatrix.PageWidth, k.templateDatamatrix.PageHeight); err != nil {
		return "", fmt.Errorf("build maroto error %w", err)
	}

	start := time.Now()

	totalItems := len(k.Cis) + len(k.Pallet)
	step := 0.0
	if totalItems > 0 {
		step = 99.0 / float64(totalItems)
	}
	for i, cis := range k.Cis {
		fnc := cis.FNC1()
		ser := cis.Serial
		pdfDocument.AddPageByTemplate(k.templateDatamatrix, fnc, ser, fmt.Sprintf("%06d", i+1))
		if ch != nil {
			ch <- step * float64(i)
		}
		// в режиме отладки берем только 10 знаков если их больше
		if k.DebugMode() {
			if i > 10 {
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

	fileName := "PDF_" + filepath.Base(model.FileCIS)
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
			return "", fmt.Errorf("save pdf %q: %w", filePdf, err)
		}
		if ch != nil {
			ch <- 100.00
		}
	}
	// запись отчета генерации
	err = pdfDocument.PdfDocumentReportSave("report.txt")
	if err != nil {
		return "", fmt.Errorf("failed to save report: %w", err)
	}
	return filePdf, nil
}
