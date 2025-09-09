package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"
)

// start начальный номер SSCC палетты
// count количество в одной палетте
func (k *Pdf) GeneratePallet(model *application.Application) error {
	indexPallet := 0
	if k.Pallet == nil {
		k.Pallet = make(map[string][]*utility.CisInfo)
	}
	for {
		k.lastSSCC = model.SsccStartNumber + indexPallet
		palet, err := utility.GenerateSSCC(k.lastSSCC, model.SsccPrefix)
		if err != nil {
			return fmt.Errorf("failed to generate SSCC: %w", err)
		}
		if _, ok := k.Pallet[palet]; ok {
			return fmt.Errorf("palet %s alredy present", palet)
		}
		// k.Pallet[palet] = k.nextRecords(indexPallet, model.PerPallet)
		k.Pallet[palet] = nextRecords(k.Cis, indexPallet, model.PerPallet)
		if len(k.Pallet[palet]) < model.PerPallet {
			// выход если сгенерировано меньше чем единиц в упаковке
			if len(k.Pallet[palet]) == 0 {
				// если пустая палета
				delete(k.Pallet, palet)
			} else {
				// если не пустая счетчик следующей палеты увеличиваем
				k.lastSSCC++
			}
			break
		}
		indexPallet++
	}
	return nil
}

// получить следующие км
// i номер группы по count штук
// если размер массива меньше count значит последний
// елси размер массива 0 значит больше нет
func (k *Pdf) nextRecords(i int, count int) (out []*utility.CisInfo) {
	lenCis := len(k.Cis)
	out = make([]*utility.CisInfo, 0)
	first := i * count // первая км в цикле 24 шт
	for i := 0; i < count; i++ {
		index := i + first
		if (index + 1) > lenCis {
			return out
		}
		out = append(out, k.Cis[index])
	}
	return out
}

// index from 0 startIndex 0
// nextRecords returns a batch of records starting from startIndex
// Returns empty slice when no more records are available
func nextRecords(arr []*utility.CisInfo, index int, count int) []*utility.CisInfo {
	startIndex := index * count
	if startIndex >= len(arr) {
		return []*utility.CisInfo{}
	}
	endIndex := startIndex + count
	// если последний индекс больше длины массива укорачиваем до размера массива
	if endIndex > len(arr) {
		endIndex = len(arr)
	}
	return arr[startIndex:endIndex]
}
