package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"
)

func (k *Pdf) GeneratePack(model *application.Application) error {
	indexPallet := 0
	if k.Pallet == nil {
		k.Pallet = make(map[string][]*utility.CisInfo)
	}
	for {
		palet := k.Kigu[indexPallet].Cis
		if _, ok := k.Pallet[palet]; ok {
			return fmt.Errorf("palet %s alredy present", palet)
		}
		k.Pallet[palet] = nextRecords(k.Cis, indexPallet, model.PerPallet)
		if len(k.Pallet[palet]) < model.PerPallet {
			// выход если сгенерировано меньше чем единиц в упаковке
			if len(k.Pallet[palet]) == 0 {
				// если пустая палета
				delete(k.Pallet, palet)
			}
			break
		}
		indexPallet++
	}
	return nil
}
