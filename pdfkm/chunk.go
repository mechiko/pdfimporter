package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"
	"slices"
)

func (k *Pdf) ChunkSplit(model *application.Application) error {
	countCisChunk := model.ChunkSize * model.PerPallet
	countCIS := 0
	countKIGU := 0
	chunksCIS := slices.Chunk(k.Cis, countCisChunk)
	k.OrderChunks = make([]string, 0)
	for chunk := range chunksCIS {
		chunkPack := &ChunkPack{
			Cis: chunk,
		}
		fileChunk := fmt.Sprintf("%06d-%06d_%s.pdf", countCIS*countCisChunk+1, ((countCIS + 1) * countCisChunk), model.FileBaseName)
		k.OrderChunks = append(k.OrderChunks, fileChunk)
		k.Chunks[fileChunk] = chunkPack
		countCIS++
	}
	chunksKIGU := slices.Chunk(k.Kigu, int(model.ChunkSize))
	for chunk := range chunksKIGU {
		fileChunk := k.OrderChunks[countKIGU]
		_, exist := k.Chunks[fileChunk]
		if !exist {
			return fmt.Errorf("ошибка поиска блока для кигу %d файл %s", countKIGU, fileChunk)
		}
		k.Chunks[fileChunk].Kigu = chunk
		countKIGU++
	}

	if countCIS != countKIGU {
		return fmt.Errorf("размер блока кигу %d не равен размеру блока км %d", countKIGU, countCIS)
	}

	return nil
}
