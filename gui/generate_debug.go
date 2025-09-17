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
		logerr("генерация пдф тест", err)
		return
	}
	model, err := GetModel()
	if err != nil {
		logerr("генерация пдф тест", err)
		return
	}
	// модель простая поэтому просто копия структуры разыменованием должно сработать
	modelStore = *model

	fileCisNormal := model.FileCIS
	if fileCisNormal != "" && utility.PathOrFileExists(fileCisNormal) {
		a.SendLog("считываем 25 КМ из файла")
		if err := pdfGenerator.ReadCIS(model); err != nil {
			logerr("генерация пдф тест ReadCSV", err)
			return
		}
		if len(pdfGenerator.Cis) < 26 {
			logerr("для примера необходимо минимум 25 марок, для печати тестовых данных не выбирайте файл КМ", err)
			return
		}
		pdfGenerator.Cis = pdfGenerator.Cis[:25]
	} else {
		a.SendLog("считываем 25 КМ из тестовых данных")
		if err := pdfGenerator.ReadCisDebug(); err != nil {
			logerr("генерация пдф тест ReadCSV", err)
			return
		}
	}

	fileKiguNormal := model.FileKIGU
	if fileKiguNormal != "" && utility.PathOrFileExists(fileKiguNormal) {
		a.SendLog("считываем 2 КИГУ из файла")
		if err := pdfGenerator.ReadKIGU(model); err != nil {
			logerr("генерация пдф тест ReadKigu", err)
			return
		}
		if len(pdfGenerator.Kigu) < 2 {
			logerr("для примера необходимо минимум 2 КИГУ, для печати тестовых данных не выбирайте файл КИГУ", err)
			return
		}
		pdfGenerator.Kigu = pdfGenerator.Kigu[:2]
	} else {
		a.SendLog("считываем 2 КИГУ из тестовых данных")
		if err := pdfGenerator.ReadKiguDebug(); err != nil {
			logerr("генерация пдф тест ReadKiguDebug", err)
			return
		}
	}

	model.FileCIS = "TEST"
	model.FileKIGU = "TEST"
	if err := pdfGenerator.GeneratePack(model); err != nil {
		logerr("gui generate", err)
		return
	}
	fileName, err := pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("генерация пдф тест", err)
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
	a.stateFinishDebug <- struct{}{}
}
