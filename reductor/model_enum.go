package reductor

import "strings"

// этот тип для ведения списка всех моделей
// просто запоминать по строке сложно вычислить ошибки в программе
// если строка вдруг окажется не такой как планировалось
type ModelType int

const (
	TrueClient ModelType = iota
	Home
	Application
	Header
	Footer
	Setup
	Index
)

// имена модели используются так же в роутинге там они выступают в качестве имен вида
// должны в роутере приводится к нижнему регистру
func (s ModelType) String() string {
	switch s {
	case Home:
		return "home"
	case TrueClient:
		return "trueclient"
	case Application:
		return "application"
	case Header:
		return "header"
	case Footer:
		return "footer"
	case Setup:
		return "setup"
	case Index:
		return "index"
	default:
		return "неизвестная"
	}
}

// строка приводится в нижний регистр потом сравнивается
func ModelTypeFromString(s string) ModelType {
	s = strings.ToLower(s)
	switch s {
	case "home":
		return Home
	case "trueclient":
		return TrueClient
	case "application":
		return Application
	case "header":
		return Header
	case "footer":
		return Footer
	case "setup":
		return Setup
	case "index":
		return Index
	}
	panic("неизвестная ModelTypeFromString")
}
