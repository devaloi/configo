package provider

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSON loads configuration from a JSON file.
type JSON struct {
	Path string
}

func NewJSON(path string) *JSON {
	return &JSON{Path: path}
}

func (p *JSON) Load() (map[string]any, error) {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, fmt.Errorf("json provider: %w", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("json provider: %w", err)
	}
	return raw, nil
}
