package pdfkm

import (
	"fmt"
	"pdfimporter/domain/models/application"
	"regexp"
	"strings"

	"pdfimporter/gs1sscc"
)

// func (k *Pdf) GenerateSSCC(i int, model *application.Application) string {
// 	code := fmt.Sprintf("%010.10s%07d", model.SsccPrefix, i)
// 	sscc := gs1sscc.Sscc(code)
// 	return "00" + sscc
// }

func (k *Pdf) GenerateSSCC(i int, model *application.Application) (string, error) {
	if model == nil {
		return "", fmt.Errorf("application model is nil")
	}
	prefix := model.SsccPrefix
	// Must be 1–10 digits
	if matched, _ := regexp.MatchString(`^\d{1,10}$`, prefix); !matched {
		return "", fmt.Errorf("invalid SsccPrefix %q: must be 1–10 digits", prefix)
	}
	// Left-zero pad or truncate to 10 digits
	switch {
	case len(prefix) < 10:
		prefix = strings.Repeat("0", 10-len(prefix)) + prefix
	case len(prefix) > 10:
		prefix = prefix[:10]
	}
	code := prefix + fmt.Sprintf("%07d", i)
	sscc := gs1sscc.Sscc(code)
	if sscc == "wrong lenght code" {
		return "", fmt.Errorf("gs1sscc.Sscc returned error for code %q", code)
	}
	return "00" + sscc, nil
}
