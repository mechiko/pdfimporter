package embeded

import (
	_ "embed"
)

//go:embed Arial-Unicode-MS.ttf
var Regular []byte

//go:embed Arial-Unicode-Bold.ttf
var Bold []byte

//go:embed Arial-Unicode-Italic.ttf
var Italic []byte

//go:embed Arial-Unicode-Bold-Italic.ttf
var BoldItalic []byte

//go:embed test.csv
var TestFile string
