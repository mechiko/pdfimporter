package gui

import (
	"fmt"
	"path/filepath"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"
)

// должна выполнятся как gorutine
func (a *GuiApp) openFile(file string) {
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
		a.Logger().Errorf("gui openFile %v", err)
		a.SendError(fmt.Sprintf("gui openFile %v", err))
		a.stateStart <- struct{}{}
		return
	}
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	model.File = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("ошибка записи модели в редуктор:", err)
		return
	}
	if model.File == "" {
		a.stateSelectedInDir <- ""
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadCSV(model); err != nil {
		logerr("ошибка загрузки файла:", err)
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ", len(pdfGenerator.Cis)))

	// устанавливаем состояни для пуск
	a.stateSelectedInDir <- filepath.Base(file)
}
