package dconfig

import (
	"fmt"
	"strconv"
	"strings"

	tk "modernc.org/tk9.0"
)

type ConfigDialogData struct {
	Ok              bool
	PrefixSSCC      string
	SsccStartNumber int
	PerPallet       int
}

type ConfigDialog struct {
	data *ConfigDialogData
	win  *tk.ToplevelWidget

	perPalet   *tk.TEntryWidget
	prefixSSCC *tk.TEntryWidget
	startSSCC  *tk.TEntryWidget

	buttonFrame  *tk.TFrameWidget
	inputFrame   *tk.TFrameWidget
	okButton     *tk.TButtonWidget
	cancelButton *tk.TButtonWidget
}

func NewConfigDialog(data *ConfigDialogData) *ConfigDialog {
	dlg := &ConfigDialog{data: data}
	dlg.win = tk.App.Toplevel()
	dlg.win.WmTitle("Config")
	// tk.WmAttributes(dlg.win, tk.Type("dialog")) // TODO
	tk.WmProtocol(dlg.win.Window, tk.WM_DELETE_WINDOW, dlg.onCancel)

	dlg.makeWidgets()
	dlg.makeLayout()
	return dlg
}

func (me *ConfigDialog) onOk() {
	me.data.Ok = true
	me.data.PrefixSSCC = fmt.Sprintf("%010s", strings.Trim(me.prefixSSCC.Textvariable(), " "))
	if start, err := strconv.ParseInt(me.startSSCC.Textvariable(), 10, 64); err == nil {
		me.data.SsccStartNumber = int(start)
	}
	if per, err := strconv.ParseInt(me.perPalet.Textvariable(), 10, 64); err == nil {
		me.data.PerPallet = int(per)
	}
	tk.Destroy(me.win)
}

func (me *ConfigDialog) onCancel() {
	tk.Destroy(me.win)
}

func (me *ConfigDialog) ShowModal() {
	me.win.Raise(tk.App)
	tk.Focus(me.win)
	// tk.Focus(me.percentSpinbox)
	tk.GrabSet(me.win)
	me.win.Center().Wait()
}
