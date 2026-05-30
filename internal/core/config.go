package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"go.yaml.in/yaml/v3"
)

type Editor struct {
	Disabled      bool              `yaml:"disabled"`
	WorkspacePath map[string]string `yaml:"workspace_path"`
}

type Config struct {
	Editors map[string]Editor `yaml:"editors"`
}

func GetConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return xdgConfigHome
	}

	return filepath.Join(os.Getenv("HOME"), ".config")
}

func GetConfigLocation() string {
	locations := map[string]string{
		"darwin": filepath.Join(GetConfigDir(), "cde.yaml"),
		"linux":  filepath.Join(GetConfigDir(), "cde.yaml"),
	}

	path, ok := locations[runtime.GOOS]
	if ok {
		return path
	}

	return filepath.Join(GetConfigDir(), "cde.yaml")
}

var defaultConfig = Config{
	Editors: map[string]Editor{
		"vscodium": {
			Disabled: false,
			WorkspacePath: map[string]string{
				"templeos": "path",
			},
		},
	},
}

var config *Config

func LoadConfig(filePath string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		config = &defaultConfig
		SaveConfig(GetConfigLocation())
		return nil
	}
	defer f.Close()

	cfg := &Config{}
	err = yaml.NewDecoder(f).Decode(cfg)
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

	err = yaml.NewEncoder(f).Encode(config)
	if err != nil {
		return fmt.Errorf("failed encoding file %s: %w", filePath, err)
	}

	return nil
}

func GetConfig() *Config {
	return config
}
