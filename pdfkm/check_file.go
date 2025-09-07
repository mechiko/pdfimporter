package pdfkm

import (
	"fmt"

	"github.com/mechiko/utility"
)

func CheckFiles(cis, kigu string, perPack int) (err error) {
	var cisArray, kiguArray []string
	if perPack <= 0 {
		return fmt.Errorf("perPack must be > 0")
	}
	if utility.PathOrFileExists(cis) {
		cisArray, err = ReadTextStringArrayFirstColon(cis)
		if err != nil {
			return fmt.Errorf("read file %s error %w", cis, err)
		}
	}
	if utility.PathOrFileExists(kigu) {
		kiguArray, err = ReadTextStringArrayFirstColon(cis)
		if err != nil {
			return fmt.Errorf("read file %s error %w", kigu, err)
		}
	}
	remainder := 0
	numberPacks := 0
	if cis != "" {
		remainder = len(cisArray) % perPack
		if remainder != 0 {
			return fmt.Errorf("количество КМ %d не кратно упаковке %d остается %d", len(cisArray), perPack, remainder)
		}
		numberPacks = len(cisArray) / perPack
		if numberPacks == 0 {
			return fmt.Errorf("количество упаковок 0")
		}
	}
	if kigu != "" {
		if cis != "" {
			if numberPacks != len(kiguArray) {
				return fmt.Errorf("в файле КИГУ: найдено %d, а необходимо %d", len(kiguArray), numberPacks)
			}
		}
	}
	return nil
}

func CheckFilesBoth(cis, kigu string, perPack int) (err error) {
	var cisArray, kiguArray []string
	if perPack <= 0 {
		return fmt.Errorf("perPack must be > 0")
	}
	if utility.PathOrFileExists(cis) {
		cisArray, err = ReadTextStringArrayFirstColon(cis)
		if err != nil {
			return fmt.Errorf("read file %s error %w", cis, err)
		}
	} else {
		return fmt.Errorf("file cis %s not found", cis)
	}
	if utility.PathOrFileExists(kigu) {
		kiguArray, err = ReadTextStringArrayFirstColon(cis)
		if err != nil {
			return fmt.Errorf("read file %s error %w", kigu, err)
		}
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
	if kigu != "" {
		if numberPacks != len(kiguArray) {
			return fmt.Errorf("в файле КИГУ: найдено %d, а необходимо %d", len(kiguArray), numberPacks)
		}
	}
	return nil
}
