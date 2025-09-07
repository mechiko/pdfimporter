package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsNew = []struct {
	name string
	err  bool
	path string
}{
	// the table itself
	{"wrong path", true, "assets"},
	{"exists path", false, "../cmd/assets"},
}

func TestNew(t *testing.T) {
	for _, tt := range testsNew {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.path)
			if tt.err {
				assert.NotNil(t, err, "ожидаем ошибку")
			} else {
				assert.NoError(t, err, "")
			}
		})
	}

}

var testsTemplates = []struct {
	name  string
	err   bool
	path  string
	count int
}{
	// the table itself
	{"exists path", false, "../cmd/assets", 2},
}

func TestTemplates(t *testing.T) {
	for _, tt := range testsTemplates {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			asts, err := New(tt.path)
			if err != nil {
				t.Errorf("%s", err.Error())
			} else {
				tmpl, err := asts.Templates()
				if tt.err {
					assert.NotNil(t, err, "")
				}
				assert.Equal(t, tt.count, len(tmpl))
			}
		})
	}

}

// 80x60 Datamatrix
func TestTemplate(t *testing.T) {
	asts, err := New("../cmd/assets")
	if err != nil {
		t.Errorf("%s", err.Error())
	} else {
		tmpl, err := asts.Template("80x60 Datamatrix")
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
		assert.Equal(t, "80x60 Datamatrix", tmpl.Name)
	}
}
