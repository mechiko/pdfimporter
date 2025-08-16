package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pdfimporter/app"
	"pdfimporter/checkdbg"
	"pdfimporter/config"
	"pdfimporter/domain/models/application"
	"pdfimporter/gui"
	"pdfimporter/pdfkm"
	"pdfimporter/reductor"
	"pdfimporter/zaplog"

	"github.com/mechiko/utility"
)

var fileExe string
var dir string

// если local true то папка создается локально
var local = flag.Bool("local", false, "")
var file = flag.String("file", "", "file to parse xlsx")

func init() {
	flag.Parse()
	fileExe = os.Args[0]
	var err error
	dir, err = filepath.Abs(filepath.Dir(fileExe))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
		os.Exit(1)
	}
}

func errMessageExit(title string, errDescription string) {
	utility.MessageBox(title, errDescription)
	os.Exit(-1)
}

func main() {
	cfg, err := config.New("", !*local)
	if err != nil {
		errMessageExit("ошибка конфигурации", err.Error())
	}

	var logsOutConfig = map[string][]string{
		"logger":   {"stdout", filepath.Join(cfg.LogPath(), config.Name)},
		"reductor": {filepath.Join(cfg.LogPath(), "reductor")},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit("ошибка создания логера", err.Error())
	}
	defer zl.Shutdown()

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit("ошибка получения логера", err.Error())
	}
	loger := lg.Sugar()
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)
	if cfg.Warning() != "" {
		loger.Infof("pkg:config warning %s", cfg.Warning())
	}

	errProcessExit := func(title string, errDescription string) {
		loger.Errorf("%s %s", title, errDescription)
		errMessageExit(title, errDescription)
	}
	// создаем приложение с опциями из конфига и логером основным
	app := app.New(cfg, loger, dir)
	// инициализируем пути необходимые приложению
	app.CreatePath()

	// создаем редуктор для хранения моделей приложения
	reductorLogger, err := zl.GetLogger("reductor")
	if err != nil {
		errProcessExit("Ошибка получения логера для редуктора", err.Error())
	}
	if err := reductor.New(reductorLogger.Sugar()); err != nil {
		errProcessExit("Ошибка создания редуктора", err.Error())
	}

	appModel, err := application.New(app)
	appModel.File = *file
	if err != nil {
		errProcessExit("Ошибка получения логера для редуктора", err.Error())
	}
	if err := reductor.Instance().SetModel(appModel, false); err != nil {
		errProcessExit("Ошибка редуктора", err.Error())
	}
	// тесты
	if err := checkdbg.NewChecks(app).Run(); err != nil {
		loger.Errorf("check error %v", err)
		errProcessExit("Check failed", err.Error())
	}

	// GUI
	k, err := pdfkm.New(app)
	if err != nil {
		errProcessExit("ошибка создание генератора пдф", err.Error())
	}
	guiApp, err := gui.New(k, app)
	if err != nil {
		errProcessExit("создание gui с ошибкой ", err.Error())
	}
	guiApp.Run()
}
