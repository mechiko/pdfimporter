package pdfkm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"pdfimporter/domain/models/application"
	"pdfimporter/embeded"
	"strings"

	"github.com/mechiko/utility"
)

func (k *Pdf) ReadCIS(model *application.Application) (err error) {
	// application.Application
	arr, err := ReadTextStringArrayFirstColon(model.FileCIS)
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

func (k *Pdf) ReadKIGU(model *application.Application) (err error) {
	// application.Application
	arr, err := ReadTextStringArrayFirstColon(model.FileKIGU)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for i, cis := range arr {
		item, err := utility.ParseCisInfo(cis)
		if err != nil {
			return fmt.Errorf("строка %d [%s] %w", i+1, cis, err)
		}
		k.Kigu = append(k.Kigu, item)
	}
	return nil
}

func (k *Pdf) ReadCisDebug() (err error) {
	// application.Application
	arr, err := readEmbeded(strings.NewReader(embeded.TestCisFile))
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

func (k *Pdf) ReadKiguDebug() (err error) {
	// application.Application
	arr, err := readEmbeded(strings.NewReader(embeded.TestKiguFile))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for i, cis := range arr {
		item, err := utility.ParseCisInfo(cis)
		if err != nil {
			return fmt.Errorf("строка %d [%s] %w", i+1, cis, err)
		}
		k.Kigu = append(k.Kigu, item)
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

func ReadTextStringArrayFirstColon(filePath string) (mp []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file %s: %w", filePath, err)
	}
	defer func() {
		if errFile := f.Close(); errFile != nil {
			// Go 1.20+: joins parse error (if any) with close error
			err = errors.Join(err, fmt.Errorf("close %s: %w", filePath, errFile))
		}
	}()

	arr := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		txtTab := strings.Split(txt, "\t")
		if len(txtTab) == 1 {
			arr = append(arr, txt)
			continue
		}
		arr = append(arr, txtTab[0])
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to parse file TXT for %s: %w", filePath, err)
	}
	return arr, nil
}
