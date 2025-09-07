package dconfig

import (
	"fmt"
	"pdfimporter/assets"

	tk "modernc.org/tk9.0"
)

func (me *ConfigDialog) makeWidgets() {
	me.makeInputs()
	me.makeButtons()
}

func (me *ConfigDialog) makeInputs() {
	me.inputFrame = me.win.TFrame()
	me.perPalet = me.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", me.data.PerPallet)))
	me.startSSCC = me.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", me.data.SsccStartNumber)))
	me.prefixSSCC = me.inputFrame.TEntry(tk.Textvariable(me.data.PrefixSSCC))
	asts, _ := assets.New("assets")
	tmplts, _ := asts.Templates()
	if me.data.MarkTemplate == "" {
		me.datamatrixCombo = me.inputFrame.TCombobox(tk.State("readonly"), tk.Textvariable("выбери шаблон"), tk.Values(tmplts))
	} else {
		me.datamatrixCombo = me.inputFrame.TCombobox(tk.State("readonly"), tk.Textvariable(me.data.MarkTemplate), tk.Values(tmplts))
	}
	if me.data.PackTemplate == "" {
		me.barCombo = me.inputFrame.TCombobox(tk.State("readonly"), tk.Textvariable("выбери шаблон"), tk.Values(tmplts))
	} else {
		me.barCombo = me.inputFrame.TCombobox(tk.State("readonly"), tk.Textvariable(me.data.MarkTemplate), tk.Values(tmplts))
	}
	me.barCombo = me.inputFrame.TCombobox(tk.State("readonly"), tk.Textvariable("выбери шаблон"), tk.Values(tmplts))
}

func (me *ConfigDialog) makeButtons() {
	me.buttonFrame = me.win.TFrame()
	me.okButton = me.buttonFrame.TButton(tk.Txt("OK"),
		tk.Command(me.onOk))
	me.cancelButton = me.buttonFrame.TButton(tk.Txt("Cancel"),
		tk.Command(me.onCancel))
}
