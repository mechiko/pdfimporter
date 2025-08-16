package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"

	"pdfimporter/gs1sscc"
)

func (k *Pdf) GenerateSSCC(i int, model *application.Application) string {
	code := fmt.Sprintf("%010.10s%07d", model.SsccPrefix, i)
	sscc := gs1sscc.Sscc(code)
	return "00" + sscc
}
