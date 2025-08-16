package config

import (
	"fmt"
	"os"
)

// тут функции надстройки возможно и не так полезные, хотя для записи значений в конфиг используются
// Viper пишет только то что прочитал изначально из конфига, если меняем что то
// то надо усатнавливать через Set а не просто в Configuration
func (c *Config) Save() error {
	if _, err := os.Stat(c.configFileName); os.IsNotExist(err) {
		if err := c.SafeWriteConfig(); err != nil {
			return fmt.Errorf("%s %w", modError, err)
		}
	} else {
		if err := c.WriteConfig(); err != nil {
			return fmt.Errorf("%s %w", modError, err)
		}
	}
	return nil
}

func (c *Config) SaveAs(fn string) error {
	err := c.WriteConfigAs(fn)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	return nil
}

func (c *Config) SaveSafe() error {
	err := c.SafeWriteConfig()
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	return nil
}

func (c *Config) GetKeyString(name string) string {
	return c.GetString(name)
}

func (c *Config) GetByName(name string) interface{} {
	return c.Get(name)
}

// записываем ключ и его значение, затем обновляем структуру Config этими значениями
// в файл не сохраняем
func (c *Config) SetInConfig(key string, value interface{}) error {
	c.Set(key, value)
	if err := c.Unmarshal(c.configuration); err != nil {
		return fmt.Errorf("Viper.Unmarshal(c.Configuration) %w", err)
	}
	return nil
}

func (c *Config) Warning() string {
	return c.warning
}

func (c *Config) DatabaseByKey(key string) *DatabaseConfiguration {
	return &DatabaseConfiguration{
		Driver: c.GetKeyString(fmt.Sprintf("%s.driver", key)),
		File:   c.GetKeyString(fmt.Sprintf("%s.file", key)),
		DbName: c.GetKeyString(fmt.Sprintf("%s.dbname", key)),
		User:   c.GetKeyString(fmt.Sprintf("%s.user", key)),
		Pass:   c.GetKeyString(fmt.Sprintf("%s.pass", key)),
		Host:   c.GetKeyString(fmt.Sprintf("%s.host", key)),
		Port:   c.GetKeyString(fmt.Sprintf("%s.port", key)),
	}

}
