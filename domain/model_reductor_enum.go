package domain

type Model string

const (
	Application  Model = "application"
	TrueClient   Model = "trueclient"
	StatusBar    Model = "statusbar"
	ZnakAgregate Model = "znakagregate"
)

func IsValidModel(s string) bool {
	switch Model(s) {
	case Application, TrueClient, StatusBar, ZnakAgregate:
		return true
	default:
		return false
	}
}
