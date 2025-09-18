package application

import (
	"fmt"
	"path/filepath"
	"pdfimporter/domain"
)

type Application struct {
	model   domain.Model
	Title   string
	Output  string
	Debug   bool
	License string

	FileBaseName    string
	FileBasePath    string
	FileCIS         string
	FileKIGU        string
	MarkTemplate    string
	PackTemplate    string
	SsccPrefix      string
	SsccStartNumber int
	PerPallet       int
	Party           string
	ChunkSize       int
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
	err = app.SetOptions("marktemplate", a.MarkTemplate)
	if err != nil {
		return fmt.Errorf("model:application save marktemplate to store error %w", err)
	}
	err = app.SetOptions("packtemplate", a.PackTemplate)
	if err != nil {
		return fmt.Errorf("model:application save packtemplate to store error %w", err)
	}
	if err := app.SaveAllOptions(); err != nil {
		return fmt.Errorf("model:application sync to store error %w", err)
	}
	err = app.SetOptions("party", a.Party)
	if err != nil {
		return fmt.Errorf("model:application save party to store error %w", err)
	}
	err = app.SetOptions("chunksize", a.ChunkSize)
	if err != nil {
		return fmt.Errorf("model:application save chunksize to store error %w", err)
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
	a.MarkTemplate = opts.MarkTemplate
	a.PackTemplate = opts.PackTemplate
	a.Party = opts.Party
	a.ChunkSize = opts.ChunkSize
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

func (a *Application) SetFileBase(file string) {
	if file == "" {
		file = "маркировка.pdf"
	}
	a.FileBasePath = filepath.Dir(file)
	filename := filepath.Base(file)
	extension := filepath.Ext(file)
	a.FileBaseName = filename[:len(filename)-len(extension)]
}
