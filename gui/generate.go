package gui

import (
	"fmt"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"

	"github.com/mechiko/utility"
)

// кнопка Пуск
// запускать в отдельном поток от tk9
func (a *GuiApp) generate() {
	logerr := func(s string, err error) {
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	defer func() {
		a.stateIsProcess <- false
	}()
	a.stateIsProcess <- true

	a.logClear <- struct{}{}
	a.SendLog("обрабатываем файлы...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	model, err := GetModel()
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadCSV(model); err != nil {
		logerr("ошибка загрузки файла:", err)
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ", len(pdfGenerator.Cis)))
	if err := pdfGenerator.GeneratePallet(model); err != nil {
		logerr("gui generate debug", err)
		return
	}
	fileName, err := pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	a.Options().SsccStartNumber = pdfGenerator.LastSSCC()
	if err := a.SaveOptions("ssccstartnumber", pdfGenerator.LastSSCC()); err != nil {
		logerr("gui generate debug", err)
		return
	}
	modelFinal, err := GetModel()
	if err != nil {
		logerr("gui generate debug", err)
		return
	}
	modelFinal.SsccStartNumber = pdfGenerator.LastSSCC()
	if err := reductor.Instance().SetModel(modelFinal, false); err != nil {
		logerr("gui generate debug", err)
		return
	}
	a.SendLog(fileName)
	if csvName, err := pdfGenerator.PaletSave("pallets"); err != nil {
		logerr("gui save palet csv", err)
		return
	} else {
		a.SendLog(csvName)
	}
	if a.DebugMode() {
		utility.OpenFileInShell(fileName)
	}
	// по завершению обработки в БД кнопка Пуск запрещена
	a.stateFinish <- struct{}{}
}
