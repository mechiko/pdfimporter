package embeded

import (
	_ "embed"
)

//go:embed RobotoCondensed-Regular.ttf
var Regular []byte

//go:embed RobotoCondensed-Bold.ttf
var Bold []byte

//go:embed RobotoCondensed-Italic.ttf
var Italic []byte

//go:embed RobotoCondensed-BoldItalic.ttf
var BoldItalic []byte

//go:embed cis.csv
var TestCisFile string

//go:embed kigu.csv
var TestKiguFile string
