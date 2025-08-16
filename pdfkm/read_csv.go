package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"
)

func (k *Pdf) ReadCSV(model *application.Application) (err error) {
	// application.Application
	arr, err := utility.ReadTextStringArray(model.File)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for i, cis := range arr {
		item, err := utility.ParseCisInfo(cis)
		if err != nil {
			return fmt.Errorf("строка %d [%s] %w", i+1, cis, err)
		}
		k.Cis = append(k.Cis, item)
	}
	return nil
}
