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
	for {
		k.lastSSCC = model.SsccStartNumber + indexPallet
		palet := k.GenerateSSCC(k.lastSSCC, model)
		if _, ok := k.Pallet[palet]; ok {
			return fmt.Errorf("palet %s alredy present", palet)
		}
		k.Pallet[palet] = k.nextRecords(indexPallet, model.PerPallet)
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
