package gui

import (
	"fmt"
	"path/filepath"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"
)

// должна выполнятся как gorutine
func (a *GuiApp) openFileCis(file string) {
	logerr := func(s string, err error) {
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	// очистка лога на экране
	a.stateIsProcess <- true
	defer func() {
		a.stateIsProcess <- false
	}()
	model, err := GetModel()
	if err != nil {
		logerr("gui openFile", err)
		return
	}
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	model.FileCIS = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("ошибка записи модели в редуктор:", err)
		return
	}
	if model.FileCIS == "" {
		a.stateSelectedCisFile <- ""
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("проверяем файл КМ")
	if err := pdfGenerator.ReadCIS(model); err != nil {
		logerr("ошибка загрузки файла:", err)
		return
	}
	remainder := len(pdfGenerator.Cis) % model.PerPallet
	numberPacks := len(pdfGenerator.Cis) / model.PerPallet
	if remainder != 0 {
		logerr("в файле КМ:", fmt.Errorf("количество КМ %d не кратно упаковке %d остается %d", len(pdfGenerator.Cis), model.PerPallet, remainder))
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ %d упаковок", len(pdfGenerator.Cis), numberPacks))
	// устанавливаем состояни для пуск
	a.stateSelectedCisFile <- filepath.Base(file)
}

func (a *GuiApp) openFileKigu(file string) {
	logerr := func(s string, err error) {
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	// очистка лога на экране
	a.stateIsProcess <- true
	defer func() {
		a.stateIsProcess <- false
	}()
	model, err := GetModel()
	if err != nil {
		logerr("gui openFile", err)
		return
	}
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	model.FileKIGU = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("ошибка записи модели в редуктор:", err)
		return
	}
	if model.FileKIGU == "" {
		a.stateSelectedCisFile <- ""
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("считываем файл КИГУ")
	if err := pdfGenerator.ReadKIGU(model); err != nil {
		logerr("ошибка загрузки файла:", err)
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КИГУ", len(pdfGenerator.Kigu)))

	// устанавливаем состояни для пуск
	a.stateSelectedKiguFile <- filepath.Base(file)
}
