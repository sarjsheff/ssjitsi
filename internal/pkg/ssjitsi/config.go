package ssjitsi

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config представляет основную конфигурацию приложения
type Config struct {
	HTTP string `yaml:"http"`
	Bots []Bot  `yaml:"bots"`
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
