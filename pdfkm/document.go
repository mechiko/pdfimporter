package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"
)

func (k *Pdf) Document(model *application.Application, ch chan float64) error {
	if k.templateDatamatrix == nil {
		return fmt.Errorf("Error pdfkm datamatrix template is nil ")
	}
	if k.templatePack == nil {
		return fmt.Errorf("Error pdfkm pack template is nil ")
	}
	totalItems := len(k.Cis)
	step := 0.0
	k.iChunkAll = 0
	k.iChunkKigu = 0
	k.iChunkCis = 0
	k.Pallet = make(map[string][]*utility.CisInfo)
	k.PackOrder = make([]string, 0)
	if totalItems > 0 {
		step = 99.0 / float64(totalItems)
	}
	iChunk := 0
	for _, file := range k.OrderChunks {
		chunk := k.Chunks[file]
		if err := k.DocumentChunk(model, ch, step, chunk, file); err != nil {
			return fmt.Errorf("ошибка генерации пдф блока %s %w", file, err)
		}
		// if k.DebugMode() {
		// 	if iChunk == 1 {
		// 		break
		// 	}
		// }
		iChunk++
	}
	if ch != nil {
		ch <- 100.00
	}
	return nil
}
