package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(file string) (*Config, error) {

	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	config := new(Config)
	if err := yaml.NewDecoder(reader).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
