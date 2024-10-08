package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Auth struct {
		Key string
	}

	Ssl struct {
		CertPath string
		KeyPath  string
	}

	Http struct {
		Port string
	}

	File struct {
		PortfolioPath string
	}
}

func NewConfig(path string) (*Config, error) {
	config := new(Config)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = toml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
