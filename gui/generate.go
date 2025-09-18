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
	if err := model.SyncToStore(a); err != nil {
		a.Logger().Errorf("ошибка синхронизации модели в настройки программы %s", err.Error())
		a.SendError(fmt.Sprintf("ошибка синхронизации модели в настройки программы %s", err.Error()))
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
		// if err := pdfGenerator.GeneratePack(model); err != nil {
		// 	logerr("генерация пдф: упаковка", err)
		// 	return
		// }
		// запрашиваем имя выходного файла и пути
		fileNamePdf := utility.TimeFileName("Этикетки") + ".pdf"
		fileNamePdfSelect, err := utility.DialogSaveFile(utility.Csv, fileNamePdf, ".")
		if err == nil {
			logerr("генерация пдф: выбор для сохранение файла агрегации", err)
			fileNamePdf = fileNamePdfSelect
		}
		model.SetFileBase(fileNamePdf)
		if err := reductor.Instance().SetModel(model, false); err != nil {
			a.Logger().Errorf("генерация пдф: ошибка сохранения модели %s", err.Error())
			a.SendError(fmt.Sprintf("генерация пдф: ошибка сохранения модели %s", err.Error()))
			return
		}
		// сплит на блоки по chunksize
		err = pdfGenerator.ChunkSplit(model)
		if err != nil {
			logerr("генерация пдф: сплит на блоки", err)
			if model != nil && model.FileCIS != "" {
				a.stateSelectedCisFile <- model.FileCIS
			}
			return
		}

		// здесь генерируем документ ПДФ целиком
		err = pdfGenerator.Document(model, a.progresCh)
		if err != nil {
			logerr("генерация пдф: документ ошибка", err)
			if model != nil && model.FileCIS != "" {
				a.stateSelectedCisFile <- model.FileCIS
			}
			return
		}
		a.SendLog("сгенерированы файлы:")
		for _, file := range pdfGenerator.Files() {
			a.SendLog(file)
		}
		fileNameCsv := utility.TimeFileName("agr_packs_"+model.FileBaseName) + ".csv"
		fileCsvSelect, err := utility.DialogSaveFile(utility.Csv, fileNameCsv, ".")
		if err != nil {
			logerr("генерация пдф: выбор для сохранение файла агрегации", err)
			fileCsvSelect = fileNameCsv
		}
		if csvName, err := pdfGenerator.PackSave(fileCsvSelect); err != nil {
			logerr("генерация пдф: сохранение файла агрегации", err)
			return
		} else {
			a.SendLog(csvName)
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
