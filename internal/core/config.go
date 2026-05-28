package core

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Editor struct {
	Disabled      bool
	WorkspacePath map[string]string
}

type Config struct {
	Editors map[string]Editor `toml:"editors" comment:"overrides for editors"`
}

var defaultConfig = Config{
	Editors: map[string]Editor{
		"vscodium": {Disabled: false},
	},
}

var config *Config

func LoadConfig(filePath string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		config = &defaultConfig
		return nil
	}
	defer f.Close()

	cfg := &Config{}
	err = toml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return fmt.Errorf("failed decoding file %s: %w", filePath, err)
	}

	config = cfg
	return nil
}

func SaveConfig(filePath string) (err error) {
	if config == nil {
		return fmt.Errorf("config file not found")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed creating file %s: %w", filePath, err)
	}
	defer f.Close()

	err = toml.NewEncoder(f).Encode(config)
	if err != nil {
		return fmt.Errorf("failed encoding file %s: %w", filePath, err)
	}

	return nil
}

func GetConfig() *Config {
	return config
}
