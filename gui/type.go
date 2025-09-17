package gui

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"pdfimporter/domain"
	"pdfimporter/domain/models/application"
	"pdfimporter/reductor"
	"time"

	tk "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	tick = 10 * time.Millisecond
)

type LogMsg struct {
	Error bool
	Msg   string
}

//go:embed 192.png
var ico []byte

type GuiApp struct {
	domain.Apper
	icon *tk.Img

	buttonFrame *tk.TFrameWidget
	inputFrame  *tk.TFrameWidget
	logFrame    *tk.TFrameWidget

	startButton  *tk.TButtonWidget
	exitButton   *tk.TButtonWidget
	debugButton  *tk.TButtonWidget
	configButton *tk.TButtonWidget
	logCh        chan LogMsg
	// stateFinishOpenXlsx   chan struct{}
	stateFinish           chan struct{}
	stateStart            chan struct{}
	logClear              chan struct{}
	stateSelectedCisFile  chan string
	stateSelectedKiguFile chan string
	stateIsProcess        chan bool
	stateFinishDebug      chan struct{}
	yscroll               *tk.Window
	logText               *tk.TextWidget

	// processing *processing.Processing
	// pdf     *pdfkm.Pdf
	fileLblCis  *tk.TLabelWidget
	fileBtnCis  *tk.TButtonWidget
	fileLblKigu *tk.TLabelWidget
	fileBtnKigu *tk.TButtonWidget

	progres   *tk.TProgressbarWidget
	progresCh chan float64
	isProcess bool
	party     *tk.TEntryWidget
	chunkSize *tk.TEntryWidget
}

func New(app domain.Apper) (*GuiApp, error) {
	a := &GuiApp{
		Apper: app,
		// pdf:   p,
	}
	a.logCh = make(chan LogMsg, 10)
	a.stateFinish = make(chan struct{})
	a.stateFinishDebug = make(chan struct{})
	a.stateStart = make(chan struct{})
	a.icon = tk.NewPhoto(tk.Data(ico))
	a.progresCh = make(chan float64, 100)
	a.logClear = make(chan struct{})
	a.stateSelectedCisFile = make(chan string, 2)
	a.stateSelectedKiguFile = make(chan string, 2)
	a.stateIsProcess = make(chan bool, 2)

	tk.App.IconPhoto(a.icon)
	tk.ErrorMode = tk.CollectErrors
	tk.App.WmTitle("Формирование ПДФ КМ")
	tk.WmProtocol(tk.App, "WM_DELETE_WINDOW", a.onQuitApp)
	if err := tk.ActivateTheme("azure light"); err != nil {
		a.Logger().Errorf("gui theme %s", err.Error())
	}
	tk.InitializeExtension("autoscroll")

	model, err := GetModel()
	if err != nil {
		return nil, fmt.Errorf("gui new get model %w", err)
	}
	a.makeWidgets(model)
	a.makeLayout()
	a.makeBindings()
	// start ticker only after widgets/layout are ready
	tk.NewTicker(tick, a.tick)
	if model.FileCIS != "" {
		go a.openFileCis(model.FileCIS)
	}
	return a, nil
}

func (a *GuiApp) Run() {
	tk.App.Center()
	tk.WmDeiconify(tk.App)
	tk.App.Wait()
}

func (a *GuiApp) logg(s, e string) {
	blue := "color1"
	red := "color2"
	if s != "" {
		s += "\n"
	}
	if e != "" {
		e += "\n"
	}
	a.logText.Configure(tk.State("normal"))
	a.logText.Insert(tk.END, s, blue, e, red)
	a.logText.See("end")
	a.logText.Configure(tk.State("disabled"))
}

func (a *GuiApp) onQuitApp() {
	if a.isProcess {
		a.logg("", "выход из программы ограничен, запущена обработка")
		return
	}
	tk.Destroy(tk.App)
}

func (a *GuiApp) tick() {
	select {
	case s := <-a.logCh:
		if s.Error {
			a.logg("", s.Msg)
		} else {
			a.logg(s.Msg, "")
		}
	case <-a.logClear:
		a.logText.Configure(tk.State("normal"))
		a.logText.Delete("1.0", tk.END)
		a.logText.Configure(tk.State("disabled"))
	case v := <-a.progresCh:
		a.progres.Configure(tk.Value(v))
	case <-a.stateStart:
		// состояние начала возможно уже выбран файл
		a.progres.Configure(tk.Value(0))
	case <-a.stateFinish:
		// состояние после записи заказов магазина в БД
		model, err := GetModel()
		if err != nil {
			a.Logger().Errorf("gui stateFinish getmodel error %s", err.Error())
		}
		model.FileCIS = ""
		model.FileKIGU = ""
		if err := reductor.Instance().SetModel(model, false); err != nil {
			a.Logger().Errorf("gui stateFinish setmodel error %s", err.Error())
		}
		a.fileLblCis.Configure(tk.Txt(""))
		a.fileLblKigu.Configure(tk.Txt(""))
		a.progres.Configure(tk.Value(0))
	case <-a.stateFinishDebug:
	case file := <-a.stateSelectedCisFile:
		label := ""
		if file != "" {
			label = filepath.Base(file)
		}
		a.fileLblCis.Configure(tk.Txt(label))
	case file := <-a.stateSelectedKiguFile:
		label := ""
		if file != "" {
			label = filepath.Base(file)
		}
		a.fileLblKigu.Configure(tk.Txt(label))
	case a.isProcess = <-a.stateIsProcess:
		if a.isProcess {
			a.fileBtnCis.Configure(tk.State("disabled"))
			a.fileBtnKigu.Configure(tk.State("disabled"))
			a.startButton.Configure(tk.State("disabled"))
			a.exitButton.Configure(tk.State("disabled"))
			a.debugButton.Configure(tk.State("disabled"))
		} else {
			a.fileBtnCis.Configure(tk.State("enabled"))
			a.fileBtnKigu.Configure(tk.State("enabled"))
			a.startButton.Configure(tk.State("enabled"))
			a.exitButton.Configure(tk.State("enabled"))
			a.debugButton.Configure(tk.State("enabled"))
		}
	default:
	}
}

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendError(s string) {
	msg := LogMsg{
		Error: true,
		Msg:   s,
	}
	a.Logger().Error(s)
	select {
	case a.logCh <- msg:
		// message sent
	default:
		// message dropped, log this event
		a.Logger().Warn("Failed to send error message to GUI: channel full")
	}
}

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendLog(s string) {
	msg := LogMsg{
		Error: false,
		Msg:   s,
	}
	select {
	case a.logCh <- msg:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send log message to GUI: channel full")
	}
}

func GetModel() (*application.Application, error) {
	modelReductor, err := reductor.Instance().Model(domain.Application)
	if err != nil {
		return nil, fmt.Errorf("failed to get model from reductor: %w", err)
	}
	model, ok := modelReductor.(*application.Application)
	if !ok {
		return nil, fmt.Errorf("model is not of type *application.Application")
	}
	return model, nil
}
