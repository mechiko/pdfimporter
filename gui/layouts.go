package gui

import tk "modernc.org/tk9.0"

func (a *GuiApp) makeLayout() {
	a.layoutLog()
	tk.Grid(a.logFrame, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS))
	a.layoutInputs()
	tk.Grid(a.inputFrame, tk.Row(1), tk.Column(0), tk.Pady(5), tk.Sticky(tk.WE))
	a.layoutButton()
	tk.Grid(a.buttonFrame, tk.Row(2), tk.Column(0), tk.Sticky(tk.WE))

	tk.GridColumnConfigure(tk.App, 0, tk.Weight(1))
	tk.GridRowConfigure(tk.App, 0, tk.Weight(1))
	tk.App.Configure(tk.Padx(5), tk.Pady(5))
}

func (a *GuiApp) layoutInputs() {
	tk.Grid(a.fileBtn, tk.Row(0), tk.Column(0), tk.Sticky(tk.WE))
	tk.Grid(a.fileLbl, tk.Row(0), tk.Column(1), tk.Sticky(tk.W))
	tk.Grid(a.progres, tk.Row(1), tk.Column(0), tk.Columnspan(2), tk.Sticky(tk.WE))

	tk.GridColumnConfigure(a.inputFrame, 1, tk.Weight(2))
	tk.GridRowConfigure(a.inputFrame, 0, tk.Pad(5))
	tk.GridRowConfigure(a.inputFrame, 1, tk.Pad(10))
}

func (a *GuiApp) layoutLog() {
	tk.Grid(a.logText, tk.Column(0), tk.Columnspan(2), tk.Sticky(tk.NEWS), tk.Pady("2m"))
	tk.Grid(a.yscroll, tk.Row(0), tk.Column(2), tk.Sticky(tk.NS+tk.E), tk.Pady("2m"))
	tk.GridRowConfigure(a.logFrame, 0, tk.Weight(1))
	tk.GridColumnConfigure(a.logFrame, 0, tk.Weight(1))
}

func (a *GuiApp) layoutButton() {
	tk.Grid(a.startButton, tk.Row(0), tk.Column(1), tk.Padx(5), tk.Sticky(tk.WE))
	tk.Grid(a.exitButton, tk.Row(0), tk.Column(2), tk.Sticky(tk.E))
}
