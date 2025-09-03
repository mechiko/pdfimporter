package gui

import (
	"fmt"
	"pdfimporter/reductor"

	"github.com/mechiko/utility"
)

// кнопка Пуск
// запускать в отдельном поток от tk9
func (a *GuiApp) generate() {
	defer func() {
		a.stateIsProcess <- false
	}()
	a.stateIsProcess <- true

	a.SendLog("обрабатываем файлы...")
	model, err := GetModel()
	if err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	if err := a.pdf.GeneratePallet(model); err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	fileName, err := a.pdf.Document(model, a.progresCh)
	if err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	a.Options().SsccStartNumber = a.pdf.LastSSCC()
	if err := a.SaveOptions("ssccstartnumber", a.pdf.LastSSCC()); err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	modelFinal, err := GetModel()
	if err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	modelFinal.SsccStartNumber = a.pdf.LastSSCC()
	if err := reductor.Instance().SetModel(modelFinal, false); err != nil {
		a.Logger().Errorf("gui generate %v", err)
		a.SendError(fmt.Sprintf("gui generate %v", err))
		a.stateStart <- struct{}{}
		return
	}
	a.SendLog(fileName)
	if csvName, err := a.pdf.PaletSave("pallets"); err != nil {
		a.Logger().Errorf("gui save palet csv %v", err)
		a.SendError(fmt.Sprintf("gui save palet csv %v", err))
		a.stateStart <- struct{}{}
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
