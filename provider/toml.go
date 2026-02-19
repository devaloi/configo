package provider

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/devaloi/configo"
)

// TOML loads configuration from a TOML file.
type TOML struct {
	Path string
}

func NewTOML(path string) *TOML {
	return &TOML{Path: path}
}

func (p *TOML) Load() (map[string]any, error) {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, fmt.Errorf("toml provider: %w", err)
	}
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("toml provider: %w", err)
	}
	return configo.Flatten(raw), nil
}
