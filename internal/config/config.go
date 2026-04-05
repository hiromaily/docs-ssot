// Package config loads the docsgen.yaml build configuration.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Target struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

type Config struct {
	Targets []Target `yaml:"targets"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
