package domain

import (
	"pdfimporter/config"

	"go.uber.org/zap"
)

type Apper interface {
	Options() *config.Configuration
	SaveOptions(key string, value interface{}) error
	SaveAllOptions() error
	Logger() *zap.SugaredLogger
	Pwd() string
	ConfigPath() string
	DbPath() string
	LogPath() string
	DebugMode() bool
}
