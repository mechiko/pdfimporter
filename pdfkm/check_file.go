package pdfkm

import (
	"fmt"

	"github.com/mechiko/utility"
)

// проверяем наличие файла
func CheckFile(file string) (size int, err error) {
	if file == "" {
		return 0, fmt.Errorf("file name empty")
	}
	if utility.PathOrFileExists(file) {
		cisArray, err := ReadTextStringArrayFirstColon(file)
		if err != nil {
			return 0, fmt.Errorf("read file %s error %w", file, err)
		}
		size = len(cisArray)
		if size == 0 {
			return 0, fmt.Errorf("read file KM 0")
		}
	}
	return size, nil
}

// проверяем наличие файлов, cis файл обязателен, perPack обязателен
// если есть кигу файл то расчитываем количество упаковок и остаток
func CheckBothFiles(cis, kigu string, perPack int) (err error) {
	var cisArray, kiguArray []string
	if perPack <= 0 {
		return fmt.Errorf("единиц в упаковке должно быть > 0")
	}
	if cis == "" {
		return fmt.Errorf("не указано имя файла КМ")
	}
	if utility.PathOrFileExists(cis) {
		cisArray, err = ReadTextStringArrayFirstColon(cis)
		if err != nil {
			return fmt.Errorf("чтение файла КМ %s ошибка %w", cis, err)
		}
		if len(cisArray) == 0 {
			return fmt.Errorf("в файле KM 0")
		}
	} else {
		return fmt.Errorf("файл КМ %s не найден", cis)
	}
	if utility.PathOrFileExists(kigu) {
		kiguArray, err = ReadTextStringArrayFirstColon(kigu)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла упаковок %s  %w", kigu, err)
		}
		if len(kiguArray) == 0 {
			return fmt.Errorf("упаковок в файле 0")
		}
		remainder := 0
		numberPacks := 0
		remainder = len(cisArray) % perPack
		if remainder != 0 {
			return fmt.Errorf("количество КМ %d не кратно упаковке %d остается %d", len(cisArray), perPack, remainder)
		}
		numberPacks = len(cisArray) / perPack
		if numberPacks == 0 {
			return fmt.Errorf("количество упаковок 0")
		}
		if numberPacks != len(kiguArray) {
			return fmt.Errorf("в файле КИГУ: найдено %d, а необходимо %d", len(kiguArray), numberPacks)
		}
	}
	return nil
}
