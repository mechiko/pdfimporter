package gui

import (
	"fmt"
	"path/filepath"
	"pdfimporter/reductor"
)

// должна выполнятся как gorutine
func (a *GuiApp) openFile(file string) {
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
	// сброс модели
	a.pdf.Reset()
	model.File = file
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		a.Logger().Errorf("gui openFile SetModel %v", err)
		a.SendError(fmt.Sprintf("ошибка записи модели в редуктор: %s", err.Error()))
		a.stateStart <- struct{}{}
		return
	}
	if model.File == "" {
		a.stateSelectedInDir <- ""
		return
	}
	a.logClear <- struct{}{}
	a.SendLog("считываем файл КМ")
	if err := a.pdf.ReadCSV(model); err != nil {
		a.Logger().Errorf("gui openFile ReadCSV %v", err)
		a.SendError(fmt.Sprintf("ошибка загрузки файла: %s", err.Error()))
		a.stateStart <- struct{}{}
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ", len(a.pdf.Cis)))

	// устанавливаем состояни для пуск
	a.stateSelectedInDir <- filepath.Base(file)
}
