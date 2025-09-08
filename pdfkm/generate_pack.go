package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"
)

func (k *Pdf) GeneratePack(model *application.Application) error {
	if k.Pallet == nil {
		k.Pallet = make(map[string][]*utility.CisInfo)
	}
	if model.PerPallet <= 0 {
		return fmt.Errorf("perPalet %d must be great 0", model.PerPallet)
	}
	if len(k.Kigu) == 0 {
		return fmt.Errorf("array kigu zero")
	}
	for indexPallet, kg := range k.Kigu {
		palet := kg.Cis
		if _, ok := k.Pallet[palet]; ok {
			return fmt.Errorf("palet %s already present", palet)
		}
		k.Pallet[palet] = nextRecords(k.Cis, indexPallet, model.PerPallet)
	}
	return nil
}
