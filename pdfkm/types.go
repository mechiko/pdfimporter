package pdfkm

import (
	"encoding/json"
	"fmt"
	"pdfimporter/assets"
	"pdfimporter/domain"
	"pdfimporter/pdfproc"

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
	templateDatamatrix *pdfproc.MarkTemplate
	templateBar        *pdfproc.MarkTemplate
}

func New(app domain.Apper) (p *Pdf, err error) {
	p = &Pdf{
		Apper:    app,
		warnings: make([]string, 0),
		errors:   make([]string, 0),
		Cis:      make([]*utility.CisInfo, 0),
		Kigu:     make([]*utility.CisInfo, 0),
		Sscc:     make([]string, 0),
		Pallet:   make(map[string][]*utility.CisInfo),
	}
	p.Reset()
	p.assets, err = assets.New("assets")
	if err != nil {
		return nil, fmt.Errorf("Error assets: %v", err)
	}
	tmplDatamatrixJson, err := p.assets.Json("datamatrix")
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}
	tmplBarJson, err := p.assets.Json("bar")
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}
	p.templateDatamatrix = &pdfproc.MarkTemplate{}
	err = json.Unmarshal(tmplDatamatrixJson, p.templateDatamatrix)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal file: %v", err)
	}
	p.templateBar = &pdfproc.MarkTemplate{}
	err = json.Unmarshal(tmplBarJson, p.templateBar)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshal file: %v", err)
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
