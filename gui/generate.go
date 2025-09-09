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
	defer func() {
		a.stateIsProcess <- false
	}()
	a.stateIsProcess <- true

	model, err := GetModel()
	if err != nil {
		a.Logger().Errorf("gui generate %s", err.Error())
		a.SendError(fmt.Sprintf("gui generate %s", err.Error()))
		return
	}
	// сохраняем модель по ошибке
	logerr := func(s string, err error) {
		if err := reductor.Instance().SetModel(model, false); err != nil {
			a.Logger().Errorf("gui generate setmodel %s", err.Error())
			a.SendError(fmt.Sprintf("gui generate setmodel  %s", err.Error()))
			return
		}
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	a.logClear <- struct{}{}
	tMark := fmt.Sprintf("выбран шаблон печати КМ: %s", model.MarkTemplate)
	tPack := fmt.Sprintf("выбран шаблон печати КИГУ: %s", model.PackTemplate)
	a.SendLog(tMark)
	a.SendLog(tPack)

	// проверяем файлы
	err = pdfkm.CheckBothFiles(model.FileCIS, model.FileKIGU, model.PerPallet)
	if err != nil {
		logerr("ошибка проверки файлов: ", err)
		return
	}

	a.SendLog("обрабатываем файлы...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("генерация пдф:", err)
		return
	}
	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadCIS(model); err != nil {
		model.FileCIS = ""
		logerr("ошибка загрузки файла:", err)
		return
	}
	numberPacks := len(pdfGenerator.Cis) / model.PerPallet
	a.SendLog(fmt.Sprintf("считано %d КМ %d упаковок", len(pdfGenerator.Cis), numberPacks))
	if model.FileKIGU != "" {
		a.SendLog("считываем файл КИГУ")
		if err := pdfGenerator.ReadKIGU(model); err != nil {
			model.FileKIGU = ""
			logerr("ошибка загрузки файла:", err)
			return
		}
		a.SendLog(fmt.Sprintf("считано %d КИГУ", len(pdfGenerator.Kigu)))
		if err := pdfGenerator.GeneratePack(model); err != nil {
			logerr("генерация пдф: упаковка", err)
			return
		}
		fileName, err := pdfGenerator.Document(model, a.progresCh)
		if err != nil {
			logerr("генерация пдф: документ", err)
			if model != nil && model.FileCIS != "" {
				a.stateSelectedCisFile <- model.FileCIS
			}
			return
		}
		a.SendLog(fileName)
		if csvName, err := pdfGenerator.PaletSave("pallets"); err != nil {
			logerr("генерация пдф: сохранение файла агрегации", err)
			return
		} else {
			a.SendLog(csvName)
		}
		if a.DebugMode() {
			utility.OpenFileInShell(fileName)
		}
	} else {
		fileName, err := pdfGenerator.DocumentWithoutPack(model, a.progresCh)
		if err != nil {
			logerr("генерация пдф: документ без упаковки", err)
			if model != nil && model.FileCIS != "" {
				a.stateSelectedCisFile <- model.FileCIS
			}
			return
		}
		a.SendLog(fileName)
		if a.DebugMode() {
			utility.OpenFileInShell(fileName)
		}
	}
	a.stateFinish <- struct{}{}
}
