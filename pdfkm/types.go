package pdfkm

import (
	"fmt"
	"pdfimporter/assets"
	"pdfimporter/domain"
	"pdfimporter/domain/models/application"
	"pdfimporter/reductor"

	"github.com/mechiko/utility"
)

// const startSSCC = "1462709225" // gs1 rus id zapivkom для памяти запивком

type Pdf struct {
	domain.Apper
	Cis                []*utility.CisInfo
	Kigu               []*utility.CisInfo
	Sscc               []string
	Pallet             map[string][]*utility.CisInfo
	lastSSCC           int
	warnings           []string
	errors             []string
	assets             *assets.Assets
	templateDatamatrix *domain.MarkTemplate
	templatePack       *domain.MarkTemplate
}

func New(app domain.Apper) (p *Pdf, err error) {
	p = &Pdf{
		Apper:              app,
		warnings:           make([]string, 0),
		errors:             make([]string, 0),
		Cis:                make([]*utility.CisInfo, 0),
		Kigu:               make([]*utility.CisInfo, 0),
		Sscc:               make([]string, 0),
		Pallet:             make(map[string][]*utility.CisInfo),
		templateDatamatrix: nil,
		templatePack:       nil,
	}
	p.Reset()
	p.assets, err = assets.New("assets")
	if err != nil {
		return nil, fmt.Errorf("Error assets: %w", err)
	}
	mdl, err := GetModel()
	if err != nil {
		return nil, fmt.Errorf("Error get model: %w", err)
	}
	if mdl.MarkTemplate != "" {
		p.templateDatamatrix, err = p.assets.Template(mdl.MarkTemplate)
		if err != nil {
			return nil, fmt.Errorf("Error get assets datamatrix template %s: %w", mdl.MarkTemplate, err)
		}
	} else {
		return nil, fmt.Errorf("не выбран шаблон печати в настройках")
	}
	if mdl.PackTemplate != "" {
		p.templatePack, err = p.assets.Template(mdl.PackTemplate)
		if err != nil {
			return nil, fmt.Errorf("Error get assets pack template %s: %w", mdl.PackTemplate, err)
		}
	} else {
		if mdl.FileKIGU != "" {
			return nil, fmt.Errorf("не выбран шаблон печати для этикеток упаковки по файлу %s", mdl.FileKIGU)
		}
	}
	return p, nil
}

func (k *Pdf) AddWarn(warn string) {
	k.warnings = append(k.warnings, warn)
}

func (k *Pdf) Warnings() []string {
	out := make([]string, len(k.warnings))
	copy(out, k.warnings)
	return out
}

func (k *Pdf) AddError(err string) {
	k.errors = append(k.errors, err)
}

func (k *Pdf) Errors() []string {
	out := make([]string, len(k.errors))
	copy(out, k.errors)
	return out
}

func (k *Pdf) Reset() {
	for key := range k.Pallet {
		delete(k.Pallet, key)
	}
	k.Sscc = make([]string, 0)
	k.Cis = make([]*utility.CisInfo, 0)
	k.Kigu = make([]*utility.CisInfo, 0)
	k.errors = make([]string, 0)
	k.warnings = make([]string, 0)
	k.lastSSCC = 0
}

func (k *Pdf) LastSSCC() int {
	return k.lastSSCC
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
