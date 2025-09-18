package gui

import (
	"fmt"
	"path/filepath"
	"pdfimporter/domain/models/application"

	"github.com/mechiko/utility"

	tk "modernc.org/tk9.0"
	"modernc.org/tk9.0/extensions/autoscroll"
)

func (a *GuiApp) makeWidgets(model *application.Application) {
	a.makeLog()
	a.makeInputs(model)
	a.makeButtons()
}

func (a *GuiApp) makeLog() {
	a.logFrame = tk.TFrame()
	a.yscroll = autoscroll.Autoscroll(tk.TScrollbar(tk.Command(func(e *tk.Event) { e.Yview(a.logText) })).Window)
	a.logText = a.logFrame.Text(
		// tk.Font(tk., 10),
		// tk.Font(tk.NewFont(tk.Family("Consolas"))),
		tk.Yscrollcommand(func(e *tk.Event) { e.ScrollSet(a.yscroll) }),
		tk.Setgrid(true),
		tk.Wrap(tk.WORD),
		tk.Padx("2m"),
		tk.Pady("2m"),
		tk.Height(10),
	)
	tag := fmt.Sprintf("t%v", 1)
	a.logText.TagConfigure(tag, tk.Foreground(tk.Red))
	a.logText.TagConfigure("bgstipple", tk.Background(tk.Black), tk.Borderwidth(0), tk.Bgstipple(tk.Gray12))
	a.logText.TagConfigure("big", tk.Font("helvetica", 12, "bold"))
	a.logText.TagConfigure("bold", tk.Font("helvetica", 10, "bold", "italic"))
	a.logText.TagConfigure("center", tk.Justify("center"))
	a.logText.TagConfigure("color1", tk.Foreground(tk.Blue))
	a.logText.TagConfigure("color2", tk.Foreground(tk.Red))
	a.logText.TagConfigure("margins", tk.Lmargin1("12m"), tk.Lmargin2("6m"), tk.Rmargin("10m"))
	a.logText.TagConfigure("overstrike", tk.Overstrike(1))
	a.logText.TagConfigure("raised", tk.Relief("raised"), tk.Borderwidth(1))
	a.logText.TagConfigure("right", tk.Justify("right"))
	a.logText.TagConfigure("spacing", tk.Spacing1("10p"), tk.Spacing2("2p"), tk.Lmargin1("12m"), tk.Lmargin2("6m"), tk.Rmargin("10m"))
	a.logText.TagConfigure("sub", tk.Offset("-2p"), tk.Font("helvetica", 8))
	a.logText.TagConfigure("sunken", tk.Relief("sunken"), tk.Borderwidth(1))
	a.logText.TagConfigure("super", tk.Offset("4p"), tk.Font("helvetica", 8))
	a.logText.TagConfigure("tiny", tk.Font("times", 8, "bold"))
	a.logText.TagConfigure("underline", tk.Underline(1))
	a.logText.TagConfigure("verybig", tk.Font(tk.CourierFont(), 22, "bold"))
}

func (a *GuiApp) makeInputs(model *application.Application) {
	a.inputFrame = tk.TFrame()
	fileCis := model.FileCIS
	labelCis := ""
	if fileCis != "" {
		base := filepath.Base(fileCis)
		if len(base) > 50 {
			labelCis = fmt.Sprintf("%.30s...%s", base, base[len(base)-10:])
		} else {
			labelCis = base
		}
	}
	fileKigu := model.FileKIGU
	labelKigu := ""
	if fileKigu != "" {
		base := filepath.Base(fileKigu)
		if len(base) > 50 {
			labelKigu = fmt.Sprintf("%.30s...%s", base, base[len(base)-10:])
		} else {
			labelKigu = base
		}
	}
	a.fileLblCis = a.inputFrame.TLabel(tk.Txt(labelCis))
	a.fileLblKigu = a.inputFrame.TLabel(tk.Txt(labelKigu))
	a.progres = a.inputFrame.TProgressbar()
	a.party = a.inputFrame.TEntry(tk.Textvariable(model.Party))
	a.chunkSize = a.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", model.ChunkSize)))
}

func (a *GuiApp) makeButtons() {
	a.buttonFrame = tk.TFrame()
	a.exitButton = a.buttonFrame.TExit(tk.Txt("Выход"))
	a.startButton = a.buttonFrame.TButton(tk.Txt("Пуск"), tk.State("enabled"), tk.Command(func() {
		if a.isProcess {
			return
		}
		a.debugButton.Configure(tk.State("disabled"))
		a.startButton.Configure(tk.State("disabled"))
		a.exitButton.Configure(tk.State("disabled"))
		go a.generate()
	}))
	a.debugButton = a.buttonFrame.TButton(tk.Txt("Пример"), tk.State("enabled"), tk.Command(func() {
		if a.isProcess {
			return
		}
		a.debugButton.Configure(tk.State("disabled"))
		a.startButton.Configure(tk.State("disabled"))
		a.exitButton.Configure(tk.State("disabled"))
		go a.generateDebug()
	}))
	a.configButton = a.buttonFrame.TButton(tk.Txt("Настройка"),
		tk.Command(func() {
			if a.isProcess {
				return
			}
			a.onConfig()
		}))

	a.fileBtnCis = a.inputFrame.TButton(tk.Txt("Открыть файл КМ"), tk.Command(func() {
		if a.isProcess {
			return
		}
		ff, err := utility.DialogOpenFile([]utility.FileType{utility.Csv, utility.All}, "", ".")
		if err != nil {
			a.logg("", err.Error())
		}
		go a.openFileCis(ff)
	}))
	a.fileBtnKigu = a.inputFrame.TButton(tk.Txt("Открыть файл КИГУ"), tk.Command(func() {
		if a.isProcess {
			return
		}
		ff, err := utility.DialogOpenFile([]utility.FileType{utility.Csv, utility.All}, "", ".")
		if err != nil {
			a.logg("", err.Error())
		}
		go a.openFileKigu(ff)
	}))
}
