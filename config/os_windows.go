package config

// пути для каталогов только относительные
// точка добавляется автоматом...
var (
	configPath       = "nevakod/" + Name
	dbPath           = "nevakod/" + Name
	logPath          = "nevakod/" + Name
	Supported        = true
	Windows          = true
	Linux            = false
	PosixUserUIDGUID = 1002
	PosixChownPath   = 0755
	PosixChownFile   = 0644
)
