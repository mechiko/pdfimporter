package pdfkm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsNew = []struct {
	name string
	err  bool
	file string
}{
	// the table itself
	{"both empty", true, ""},
	{"both empty", false, "../.DATA/Заказ_00000000377_Пиво светлое пастеризованное фильтрованное «Харп Лагер»_ГТИН_05000213100066_2376.csv"},
}

func TestCheck(t *testing.T) {
	for _, tt := range testsNew {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckFile(tt.file)
			if tt.err {
				assert.NotNil(t, err, "ожидаем ошибку")
			} else {
				assert.NoError(t, err, "")
			}
		})
	}

}
