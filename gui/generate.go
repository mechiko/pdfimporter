package gui

import (
	"errors"
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
	a.SendLog("обрабатываем файлы...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("gui generate", err)
		return
	}

	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadCIS(model); err != nil {
		model.FileCIS = ""
		logerr("ошибка загрузки файла:", err)
		return
	}
	if len(pdfGenerator.Cis) == 0 {
		model.FileCIS = ""
		logerr("в файле КМ:", errors.New("0 cis"))
		return
	}
	remainder := len(pdfGenerator.Cis) % model.PerPallet
	numberPacks := len(pdfGenerator.Cis) / model.PerPallet
	if remainder != 0 {
		model.FileCIS = ""
		logerr("в файле КМ:", fmt.Errorf("количество КМ %d не кратно упаковке %d остается %d", len(pdfGenerator.Cis), model.PerPallet, remainder))
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КМ %d упаковок", len(pdfGenerator.Cis), numberPacks))

	a.SendLog("считываем файл КИГУ")
	if err := pdfGenerator.ReadKIGU(model); err != nil {
		model.FileKIGU = ""
		logerr("ошибка загрузки файла:", err)
		return
	}
	if len(pdfGenerator.Kigu) == 0 {
		model.FileKIGU = ""
		logerr("в файле КИГУ:", errors.New("0 KIGU"))
		return
	}
	a.SendLog(fmt.Sprintf("считано %d КИГУ", len(pdfGenerator.Kigu)))

	if len(pdfGenerator.Kigu) != numberPacks {
		model.FileKIGU = ""
		logerr("в файле КИГУ:", fmt.Errorf("найдено %d, а необходимо %d", len(pdfGenerator.Kigu), numberPacks))
		return
	}
	if err := pdfGenerator.GeneratePallet(model); err != nil {
		logerr("gui generate debug", err)
		return
	}
	fileName, err := pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("gui generate debug", err)
		if model != nil && model.FileCIS != "" {
			a.stateSelectedCisFile <- model.FileCIS
		}
		return
	}
	a.Options().SsccStartNumber = pdfGenerator.LastSSCC()
	if err := a.SetOptions("ssccstartnumber", pdfGenerator.LastSSCC()); err != nil {
		logerr("gui generate", err)
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
