package gui

import (
	"fmt"
	"pdfimporter/domain/models/application"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"

	"github.com/mechiko/utility"
)

// кнопка Пуск
// запускать в отдельном поток от tk9
func (a *GuiApp) generateDebug() {
	modelStore := application.Application{}
	logerr := func(s string, err error) {
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	defer func() {
		// восстанавливаем модель
		if err := reductor.Instance().SetModel(&modelStore, false); err != nil {
			a.Logger().Errorf("%s", err.Error())
		}
		a.stateIsProcess <- false
	}()
	a.stateIsProcess <- true

	a.logClear <- struct{}{}
	a.SendLog("обрабатываем тест...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	model, err := GetModel()
	if err != nil {
		logerr("gui generate", err)
		return
	}
	// модель простая поэтому просто копия структуры разыменованием должно сработать
	modelStore = *model

	model.File = "TEST"
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		logerr("gui openFile SetModel", err)
		return
	}
	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadDebug(); err != nil {
		logerr("gui openFile ReadCSV", err)
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ", len(pdfGenerator.Cis)))

	if err := pdfGenerator.GeneratePallet(model); err != nil {
		logerr("gui generate", err)
		return
	}
	fileName, err := pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("gui generate", err)
		return
	}
	a.SendLog(fileName)
	utility.OpenFileInShell(fileName)
	// по завершению обработки в БД кнопка Пуск запрещена
	a.stateFinish <- struct{}{}
}
