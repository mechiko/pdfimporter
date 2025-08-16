package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	HomePath       string
	absConfigPath  string
	absDbPath      string
	absLogPath     string
	configuration  *Configuration
	configFileName string
	warning        string
}

const modError = "pkg:config"

func New(cfgName string, inUserHome bool) (cfg *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic %v", modError, r)
		}
	}()

	configName := cfgName
	if cfgName == "" {
		configName = "config"
	}

	// viperOrigin := viper.GetViper()
	viperOrigin := viper.New()
	cfg = &Config{
		Viper:    viperOrigin,
		HomePath: ".",
	}
	if inUserHome {
		cfg.HomePath = osUserHomeDir()
	}
	if err := cfg.initPaths(); err != nil {
		return nil, fmt.Errorf("%s: failed to initialize paths: %w", modError, err)
	}
	cfg.configFileName = filepath.Join(cfg.absConfigPath, configName+".toml")
	viperConfigPath := configPath
	if !strings.HasPrefix(viperConfigPath, ".") {
		viperConfigPath = "." + viperConfigPath
	}
	viperConfigPath = filepath.Join(cfg.HomePath, viperConfigPath)
	viperOrigin.SetConfigName(configName)
	viperOrigin.SetConfigType("toml")
	viperOrigin.AddConfigPath(viperConfigPath)

	if err := viperOrigin.MergeConfig(strings.NewReader(string(TomlConfig))); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	if err := viperOrigin.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// другая ошибка
			return nil, fmt.Errorf("%s %w", modError, err)
		} else {
			// конфиг файл не найден
			cfg.warning = fmt.Sprintf("config file ('%s') not found\n", cfg.configFileName)
		}
	}
	cfg.configuration = &Configuration{}
	if err := viperOrigin.Unmarshal(cfg.configuration); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	// игнорируем ошибку этот вызов для первого сохранения файла конфига когда его нет
	err = viperOrigin.SafeWriteConfig()
	if err != nil {
		cfg.warning += fmt.Sprintf("%s\n", err.Error())
	}
	return cfg, nil
}

// вроде как возвращает копию структуры через разыменование
// проверено :)
func (c *Config) Configuration() *Configuration {
	sc := *c.configuration
	return &sc
}

func osUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to HOME environment variable
		home = os.Getenv("HOME")
		if home == "" {
			// Return a safe default or handle the error appropriately
			home = "."
		}
	}
	return home
}

func (c *Config) initPaths() (err error) {
	// создаем папки по конфигурации
	if c.absConfigPath, err = c.createConfigDirectory(configPath); err != nil {
		return err
	}
	if c.absLogPath, err = c.createConfigDirectory(logPath); err != nil {
		return err
	}
	if c.absDbPath, err = c.createConfigDirectory(dbPath); err != nil {
		return err
	}
	return nil
}

// пути в конфиге задаются только относительные!
func (c *Config) createConfigDirectory(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path is absolute %s", path)
	}
	if !strings.HasPrefix(path, ".") {
		path = "." + path
	}
	fullPath := filepath.Join(c.HomePath, path)
	if err := pathCreate(fullPath); err != nil && !errors.Is(err, fs.ErrExist) {
		return "", fmt.Errorf("невозможно создать путь %s %w", fullPath, err)
	}
	return filepath.Abs(fullPath)
}

// создаст весь путь вложенных каталогов
func pathCreate(path string) error {
	if path != "" {
		if err := os.MkdirAll(path, os.ModePerm); err != nil { // создает весь путь
			// if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) ConfigPath() string {
	return c.absConfigPath
}

func (c *Config) DbPath() string {
	return c.absDbPath
}

func (c *Config) LogPath() string {
	return c.absLogPath
}
