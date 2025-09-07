package gui

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

}
