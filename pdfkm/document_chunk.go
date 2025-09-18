package pdfkm

import (
	"fmt"
	"path/filepath"
	"pdfimporter/domain/models/application"
	"pdfimporter/pdfproc"
	"slices"
)

func (k *Pdf) DocumentChunk(model *application.Application, ch chan float64, step float64, chunk *ChunkPack, file string) error {
	if k.templateDatamatrix == nil {
		return fmt.Errorf("Error pdfkm datamatrix template is nil ")
	}
	if k.templatePack == nil {
		return fmt.Errorf("Error pdfkm pack template is nil ")
	}
	pdfDocument, err := pdfproc.New(k, k.assets)
	if err != nil {
		return fmt.Errorf("Error create pdfproc: %v", err)
	}
	if err := pdfDocument.BuildMaroto(k.templateDatamatrix.PageWidth, k.templateDatamatrix.PageHeight); err != nil {
		return fmt.Errorf("build maroto error %w", err)
	}
	cises := chunk.Cis
	kigus := chunk.Kigu
	iKigu := 0
	party := fmt.Sprintf("%.2s", model.Party)
	packs := slices.Chunk(cises, model.PerPallet)
	for ciss := range packs {
		// генерируем км упаковки
		for _, cis := range ciss {
			k.iChunkAll++
			k.iChunkCis++
			// генерируем КМ
			fnc := cis.FNC1()
			// ser := cis.Serial
			if err := pdfDocument.AddPageByTemplate(k.templateDatamatrix, fnc, party, fmt.Sprintf("%06d", k.iChunkCis)); err != nil {
				return fmt.Errorf("add datamatrix KM in page (idx %d): %w", k.iChunkAll, err)
			}
			if ch != nil {
				ch <- step * float64(k.iChunkAll)
			}
		}
		// генерируем KIGU
		if iKigu >= len(kigus) {
			return fmt.Errorf("ошибка индекса iKigu in page %d", iKigu)
		}
		kigu := kigus[iKigu]
		if _, exist := k.Pallet[kigu.Cis]; exist {
			return fmt.Errorf("ошибка упаковка %s уже обработана ранее", kigu.Cis)
		}
		k.PackOrder = append(k.PackOrder, kigu.Cis)
		k.Pallet[kigu.Cis] = ciss
		fnc := kigu.FNC1()
		// ser := kigu.Serial
		if err := pdfDocument.AddPageByTemplate(k.templatePack, fnc, party, fmt.Sprintf("%06d", k.iChunkKigu+1)); err != nil {
			return fmt.Errorf("add pack in page (idx %d): %w", k.iChunkKigu+1, err)
		}
		iKigu++
		if ch != nil {
			ch <- step * float64(k.iChunkAll)
		}
		k.iChunkKigu++
	}
	// for _, cis := range cises {
	// 	k.iChunkAll++
	// 	k.iChunkCis++
	// 	// генерируем КМ
	// 	fnc := cis.FNC1()
	// 	// ser := cis.Serial
	// 	if err := pdfDocument.AddPageByTemplate(k.templateDatamatrix, fnc, party, fmt.Sprintf("%06d", k.iChunkCis)); err != nil {
	// 		return fmt.Errorf("add datamatrix KM in page (idx %d): %w", k.iChunkAll, err)
	// 	}
	// 	if ch != nil {
	// 		ch <- step * float64(k.iChunkAll)
	// 	}
	// 	lastInPack := (k.iChunkCis % model.PerPallet) == 0
	// 	if lastInPack {
	// 		// генерируем KIGU
	// 		if iKigu >= len(kigus) {
	// 			return fmt.Errorf("ошибка индекса iKigu in page %d", iKigu)
	// 		}
	// 		kigu := kigus[iKigu]
	// 		fnc := kigu.FNC1()
	// 		// ser := kigu.Serial
	// 		if err := pdfDocument.AddPageByTemplate(k.templatePack, fnc, party, fmt.Sprintf("%06d", k.iChunkKigu+1)); err != nil {
	// 			return fmt.Errorf("add pack in page (idx %d): %w", k.iChunkKigu+1, err)
	// 		}
	// 		iKigu++
	// 		if ch != nil {
	// 			ch <- step * float64(k.iChunkAll)
	// 		}
	// 		k.iChunkKigu++
	// 	}
	// }
	err = pdfDocument.DocumentGenerate()
	if err != nil {
		return fmt.Errorf("генерация пдф блока ошибка %w", err)
	}
	err = pdfDocument.PdfDocumentSave(filepath.Join(model.FileBasePath, file))
	if err != nil {
		return fmt.Errorf("генерация пдф блока save error %q: %w", file, err)
	}
	return nil
}
