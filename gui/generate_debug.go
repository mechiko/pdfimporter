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

	fileCisNormal := model.FileCIS
	a.SendLog("считываем файл КМ")
	if fileCisNormal != "" && utility.PathOrFileExists(fileCisNormal) {
		if err := pdfGenerator.ReadCIS(model); err != nil {
			logerr("gui openFile ReadCSV", err)
			return
		}
		if len(pdfGenerator.Cis) > 25 {
			pdfGenerator.Cis = pdfGenerator.Cis[:25]
		}
	} else {
		if err := pdfGenerator.ReadCisDebug(); err != nil {
			logerr("gui ReadDebug debug ReadCSV", err)
			return
		}
	}
	a.SendLog(fmt.Sprintf("считано %d КМ", len(pdfGenerator.Cis)))

	fileKiguNormal := model.FileCIS
	a.SendLog("считываем файл КИГУ")
	if fileKiguNormal != "" && utility.PathOrFileExists(fileKiguNormal) {
		if err := pdfGenerator.ReadKIGU(model); err != nil {
			logerr("gui openFile ReadKigu", err)
			return
		}
		if len(pdfGenerator.Kigu) > 1 {
			pdfGenerator.Kigu = pdfGenerator.Kigu[:2]
		}
	} else {
		if err := pdfGenerator.ReadKiguDebug(); err != nil {
			logerr("gui ReadDebug debug ReadKiguDebug", err)
			return
		}
	}
	a.SendLog(fmt.Sprintf("считано %d КИГУ", len(pdfGenerator.Kigu)))

	model.FileCIS = "TEST"
	model.FileKIGU = "TEST"
	if err := pdfGenerator.GeneratePallet(model); err != nil {
		logerr("gui generate", err)
		return
	}
	fileName, err := pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("gui generate", err)
		if modelStore.FileCIS != "" {
			a.stateSelectedCisFile <- modelStore.FileCIS
		}
		return
	}
	a.SendLog(fileName)
	utility.OpenFileInShell(fileName)
	if modelStore.FileCIS != "" {
		a.stateSelectedCisFile <- modelStore.FileCIS
		return
	}
	a.stateFinish <- struct{}{}
}
