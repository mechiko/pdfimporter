package dconfig

import (
	tk "modernc.org/tk9.0"
)

func (me *ConfigDialog) makeBindings() {
	tk.Bind(me.datamatrixCombo, "<<ComboboxSelected>>", tk.Command(func() {
		me.data.MarkTemplate = me.datamatrixCombo.Textvariable()
	}))
	tk.Bind(me.packCombo, "<<ComboboxSelected>>", tk.Command(func() {
		me.data.PackTemplate = me.packCombo.Textvariable()
	}))

}
