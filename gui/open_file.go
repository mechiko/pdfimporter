package gui

import (
	"fmt"
	"path/filepath"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"
)

// должна выполнятся как gorutine
func (a *GuiApp) openFileCis(file string) {
	a.cis = file
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
	model.FileCIS = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("ошибка записи модели в редуктор:", err)
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("проверяем файл КМ")
	err = pdfkm.CheckFiles(a.cis, a.kigu, model.PerPallet)
	if err != nil {
		logerr("ошибка проверки файлов: ", err)
		return
	}
	// устанавливаем состояни для пуск
	a.stateSelectedCisFile <- filepath.Base(file)
}

func (a *GuiApp) openFileKigu(file string) {
	a.kigu = file
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
	model.FileKIGU = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("ошибка записи модели в редуктор:", err)
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("проверяем файл КИГУ")
	err = pdfkm.CheckFiles(a.cis, a.kigu, model.PerPallet)
	if err != nil {
		logerr("ошибка проверки файлов: ", err)
		return
	}
	// устанавливаем состояни для пуск
	a.stateSelectedKiguFile <- filepath.Base(file)
}
