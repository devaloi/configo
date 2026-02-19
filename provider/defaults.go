package provider

import (
	"github.com/devaloi/configo"
)

// Defaults provides static default values.
type Defaults struct {
	Values map[string]any
}

func NewDefaults(values map[string]any) *Defaults {
	return &Defaults{Values: values}
}

func (p *Defaults) Load() (map[string]any, error) {
	return configo.Flatten(p.Values), nil
}
