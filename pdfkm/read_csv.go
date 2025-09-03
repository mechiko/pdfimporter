package pdfkm

import (
	"bufio"
	"fmt"
	"io"
	"pdfimporter/domain/models/application"
	"pdfimporter/embeded"
	"strings"

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

func (k *Pdf) ReadDebug() (err error) {
	// application.Application
	arr, err := readEmbeded(strings.NewReader(embeded.TestFile))
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

func readEmbeded(file io.Reader) (mp []string, err error) {
	arr := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to parse file TXT: %w", err)
	}
	return arr, nil
}
