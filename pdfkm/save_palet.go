package pdfkm

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mechiko/utility"
)

func (k *Pdf) PackSave(fileName string) (string, error) {
	csvPalet := make([][]string, 0)
	for _, palet := range k.PackOrder {
		cises := k.Pallet[palet]
		for _, cis := range cises {
			csvPalet = append(csvPalet, []string{cis.Cis, palet})
		}
	}
	if err := saveTxtCustom(fileName, csvPalet); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return fileName, nil
}

func (k *Pdf) PaletSave(fileName string) (string, error) {
	fileName = utility.TimeFileName(fileName)
	if !strings.HasSuffix(fileName, ".csv") {
		fileName = fmt.Sprintf("%s.csv", fileName)
	}
	palets := make([]string, 0, len(k.Pallet))
	for k2 := range k.Pallet {
		palets = append(palets, k2)
	}
	slices.Sort(palets)

	csvPalet := make([][]string, 0)
	for _, palet := range palets {
		cises := k.Pallet[palet]
		for _, cis := range cises {
			csvPalet = append(csvPalet, []string{cis.Cis, palet})
		}
	}
	if err := saveTxtCustom(fileName, csvPalet); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return fileName, nil
}
