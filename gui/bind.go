package gui

import (
	"pdfimporter/domain"
	"pdfimporter/domain/models/application"
	"pdfimporter/reductor"
	"strconv"

	tk "modernc.org/tk9.0"
)

func (a *GuiApp) makeBindings() {
	// tk.Bind(tk.App, "<Escape>", tk.Command(a.onQuitApp))
	// tk.Bind(tk.App, "<<ComboboxSelected>>", tk.Command(func() {
	// 	model, err := GetModel()
	// 	if err != nil {
	// 		a.Logger().Errorf("gui new get model %w", err)
	// 		return
	// 	}
	// 	model.Magazin = a.magazinCombo.Textvariable()
	// 	// a.magazinCombo.Configure(tk.Textvariable(model.Magazin))
	// 	if _, ok := model.Reestr[model.Magazin]; ok {
	// 		a.SendLog(fmt.Sprintf("выбран магазин %s", model.Magazin))
	// 		a.startButton.Configure(tk.State("enabled"))
	// 		reductor.Instance().SetModel(model, false)
	// 	} else {
	// 		a.startButton.Configure(tk.State("disabled"))
	// 	}
	// }))
	tk.Bind(a.party, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		party := a.party.Textvariable()
		mdl, _ := reductor.Instance().Model(domain.Application)
		model, ok := mdl.(*application.Application)
		if !ok {
			a.Logger().Errorf("bad type model aplication")
			return
		}
		model.Party = party
		reductor.Instance().SetModel(model, false)
	}))
	tk.Bind(a.chunkSize, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		chunkSize, err := strconv.ParseInt(a.chunkSize.Textvariable(), 10, 64)
		if err != nil {
			a.chunkSize.Configure(tk.Textvariable(""))
		}
		mdl, _ := reductor.Instance().Model(domain.Application)
		model, ok := mdl.(*application.Application)
		if !ok {
			a.Logger().Errorf("bad type model aplication")
		}
		model.ChunkSize = chunkSize
		reductor.Instance().SetModel(model, false)
	}))

}
