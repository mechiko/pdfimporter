package gui

import (
	"pdfimporter/gui/dconfig"
	"pdfimporter/reductor"
)

func (a *GuiApp) onConfig() {
	model, _ := GetModel()
	data := dconfig.ConfigDialogData{
		PrefixSSCC:      model.SsccPrefix,
		PerPallet:       model.PerPallet,
		SsccStartNumber: model.SsccStartNumber,
		MarkTemplate:    model.MarkTemplate,
		PackTemplate:    model.PackTemplate,
	}
	dlg := dconfig.NewConfigDialog(&data)
	dlg.ShowModal()
	if data.Ok {
		model.PerPallet = data.PerPallet
		model.SsccPrefix = data.PrefixSSCC
		model.SsccStartNumber = data.SsccStartNumber
		model.MarkTemplate = data.MarkTemplate
		model.PackTemplate = data.PackTemplate
		if err := model.SyncToStore(a); err != nil {
			a.Logger().Errorf("диалог onConfig синхронизация модели %v", err)
		}
		err := reductor.Instance().SetModel(model, false)
		if err != nil {
			a.Logger().Errorf("dialog onConfig set reductor error %v", err)
		}
	}
}
