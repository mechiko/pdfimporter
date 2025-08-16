package application

import (
	"fmt"
	"pdfimporter/config"
	"pdfimporter/domain"
)

type Application struct {
	model   domain.Model
	Title   string
	Output  string
	Debug   bool
	License string
	//
	Height          int64
	Width           int64
	File            string
	SsccPrefix      string `json:"ssccprefix"`
	SsccStartNumber int    `json:"ssccstartnumber"`
	PerPallet       int    `json:"perpallet"`
}

var _ domain.Modeler = (*Application)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper) (*Application, error) {
	model := &Application{
		model:  domain.Application,
		Title:  "Application Title",
		Height: 60,
		Width:  80,
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model application read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (a *Application) SyncToStore(app domain.Apper) (err error) {
	// ...
	return err
}

// читаем состояние приложения
func (a *Application) ReadState(app domain.Apper) (err error) {
	a.Debug = config.Mode == "development"
	a.SsccPrefix = app.Options().SsccPrefix
	a.SsccStartNumber = app.Options().SsccStartNumber
	a.PerPallet = app.Options().PerPallet
	return nil
}

func (a *Application) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *Application) Model() domain.Model {
	return a.model
}

func (a *Application) Reset() {
}
