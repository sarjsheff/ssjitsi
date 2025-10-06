package ssjitsi

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config представляет основную конфигурацию приложения
type Config struct {
	HTTP        string `yaml:"http"`
	WebUsername string `yaml:"web_username"` // Логин для доступа к веб-консоли
	WebPassword string `yaml:"web_password"` // Пароль для доступа к веб-консоли
	Bots        []Bot  `yaml:"bots"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
