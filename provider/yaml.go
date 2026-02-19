package provider

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YAML loads configuration from a YAML file.
type YAML struct {
	Path string
}

func NewYAML(path string) *YAML {
	return &YAML{Path: path}
}

func (p *YAML) Load() (map[string]any, error) {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, fmt.Errorf("yaml provider: %w", err)
	}
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("yaml provider: %w", err)
	}
	return raw, nil
}
