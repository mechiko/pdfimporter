package pdfkm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsNew = []struct {
	name string
	err  bool
	cis  string
	kigu string
	per  int
}{
	// the table itself
	{"both empty", false, "", "", 24},
	{"both empty", false, "../.DATA/Заказ_00000000377_Пиво светлое пастеризованное фильтрованное «Харп Лагер»_ГТИН_05000213100066_2376.csv", "", 24},
	{"both empty", false, "", "../.DATA/Заказ_00000000377_Пиво светлое пастеризованное фильтрованное «Харп Лагер»_ГТИН_05000213100066_2376.csv", 24},
	{"both empty", true, "../.DATA/Заказ_00000000377_Пиво светлое пастеризованное фильтрованное «Харп Лагер»_ГТИН_05000213100066_2376.csv", "../.DATA/Заказ_00000000377_Пиво светлое пастеризованное фильтрованное «Харп Лагер»_ГТИН_05000213100066_2376.csv", 24},
	{"both empty", true, ".DATA/Заказ_2376.csv", "", 24},
}

func TestCheck(t *testing.T) {
	for _, tt := range testsNew {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := CheckFiles(tt.cis, tt.kigu, 24)
			if tt.err {
				assert.NotNil(t, err, "ожидаем ошибку")
			} else {
				assert.NoError(t, err, "")
			}
		})
	}

}
