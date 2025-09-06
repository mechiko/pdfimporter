package app

import (
	"context"
	"fmt"

	"pdfimporter/config"
	"pdfimporter/domain"

	// "pdfimporter/repo"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// type IApp interface {
// 	Options() *config.Configuration
// 	SaveOptions(string, any) error
// 	Logger() *zap.SugaredLogger
// }

type app struct {
	ctx       context.Context
	uuid      string // идентификатор для уникальности формы
	config    *config.Config
	options   *config.Configuration // копия config.Configuration
	loger     *zap.SugaredLogger
	pwd       string
	startTime time.Time
	endTime   time.Time
	// repo      *repo.Repository
	output string
}

var _ domain.Apper = (*app)(nil)

// const modError = "app"

func New(cfg *config.Config, logger *zap.SugaredLogger, pwd string) *app {
	newApp := &app{}
	newApp.pwd = pwd
	newApp.loger = logger
	newApp.config = cfg
	newApp.options = cfg.Configuration()
	newApp.uuid = uuid.New().String()
	newApp.initDateMn()
	// newApp.options.Export = "local copy"
	// if err := newApp.SaveOptions("export", "config copy"); err != nil {
	// 	fmt.Println(err)
	// }
	return newApp
}

func (a *app) initDateMn() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		a.loger.Warnf("failed to load timezone, using UTC: %v", err)
		loc = time.UTC
	}
	t := time.Now().In(loc)
	year, month, _ := t.Date()
	a.startTime = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	a.endTime = time.Date(year, month+1, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
}

func (a *app) NowDateString() string {
	n := time.Now()
	return fmt.Sprintf("%4d.%02d.%02d %02d:%02d:%02d", n.Local().Year(), n.Local().Month(), n.Local().Day(), n.Local().Hour(), n.Local().Minute(), n.Local().Second())
}

func (a *app) StartDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.startTime.Local().Year(), a.startTime.Local().Month(), a.startTime.Local().Day())
}

func (a *app) EndDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.endTime.Local().Year(), a.endTime.Local().Month(), a.endTime.Local().Day())
}

func (a *app) SetStartDate(d time.Time) {
	a.startTime = d
}

func (a *app) SetEndDate(d time.Time) {
	a.endTime = d
}

func (a *app) StartDate() time.Time {
	return a.startTime
}

func (a *app) EndDate() time.Time {
	return a.endTime
}

// func (a *app) SetRepo(repo *repo.Repository) {
// 	a.repo = repo
// }

// func (a *app) FsrarID() string {
// 	return a.Config().Configuration().Application.Fsrarid
// }

func (a *app) SetFsrarID(id string) {
	a.Config().SetInConfig("application.fsrarid", id)
}

func (a *app) Pwd() string {
	return a.pwd
}

// func (a *app) Repo() *repo.Repository {
// 	return a.repo
// }

func (a *app) Output() string {
	return a.output
}

func (a *app) Config() *config.Config {
	return a.config
}

func (a *app) Logger() *zap.SugaredLogger {
	return a.loger
}

func (a *app) Ctx() context.Context {
	return a.ctx
}

// выдаем адрес структуры опций программы чтобы править по месту
func (a *app) Options() *config.Configuration {
	return a.options
}

// записываем ключ и его значение только в пакет config
// изменения записываются в файл конфигурации
func (a *app) SetOptions(key string, value any) error {
	a.config.SetInConfig(key, value)
	if err := a.config.Save(); err != nil {
		return fmt.Errorf("save in config error %w", err)
	}
	return nil
}

// изменения записываются в файл конфигурации
func (a *app) SaveAllOptions() error {
	if err := a.config.Save(); err != nil {
		return fmt.Errorf("save all in config error %w", err)
	}
	return nil
}

// создаем по необходимости пути программы
func (a *app) CreatePath() error {
	// создаем папку вывода если не пустое значение
	// в папке запуска программы только или если она задана абсолютным значением пути
	if a.options == nil {
		return fmt.Errorf("опции программы не инициализированы")
	}
	// if a.options.Output != "" {
	// 	if output, err := createPath(a.options.Output, ""); err != nil {
	// 		return fmt.Errorf("ошибка создания каталога %w", err)
	// 	} else {
	// 		a.options.Output = output
	// 	}
	// 	a.loger.Infof("путь output приложения %s", a.options.Output)
	// }
	return nil
}

func (a *app) ConfigPath() string {
	if a.config != nil {
		return a.config.ConfigPath()
	}
	return ""
}

func (a *app) DbPath() string {
	if a.config != nil {
		return a.config.DbPath()
	}
	return ""
}

func (a *app) LogPath() string {
	if a.config != nil {
		return a.config.LogPath()
	}
	return ""
}

func (a *app) DebugMode() bool {
	if config.Mode == "development" {
		return true
	}
	return false
}
