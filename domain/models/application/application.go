package application

import (
	"fmt"
	"pdfimporter/domain"
)

type Application struct {
	model   domain.Model
	Title   string
	Output  string
	Debug   bool
	License string

	FileCIS         string
	FileKIGU        string
	SsccPrefix      string `json:"ssccprefix"`
	SsccStartNumber int    `json:"ssccstartnumber"`
	PerPallet       int    `json:"perpallet"`
}

var _ domain.Modeler = (*Application)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper) (*Application, error) {
	model := &Application{
		model: domain.Application,
		Title: "Application Title",
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model application read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (a *Application) SyncToStore(app domain.Apper) (err error) {
	// ...
	err = app.SetOptions("ssccprefix", a.SsccPrefix)
	if err != nil {
		return fmt.Errorf("model:application save ssccprefix to store error %w", err)
	}
	err = app.SetOptions("ssccstartnumber", a.SsccStartNumber)
	if err != nil {
		return fmt.Errorf("model:application save ssccstartnumber to store error %w", err)
	}
	err = app.SetOptions("perpallet", a.PerPallet)
	if err != nil {
		return fmt.Errorf("model:application save perpallet to store error %w", err)
	}
	if err := app.SaveAllOptions(); err != nil {
		return fmt.Errorf("model:application sync to store error %w", err)
	}
	return err
}

// читаем состояние приложения
func (a *Application) ReadState(app domain.Apper) (err error) {
	opts := app.Options()
	if opts == nil {
		return fmt.Errorf("nil options from app")
	}
	a.Debug = app.DebugMode()
	a.SsccPrefix = opts.SsccPrefix
	a.SsccStartNumber = opts.SsccStartNumber
	a.PerPallet = opts.PerPallet
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
