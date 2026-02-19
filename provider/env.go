package provider

import (
	"os"
	"strings"
)

// Env loads configuration from environment variables with prefix filtering.
type Env struct {
	Prefix string
}

func NewEnv(prefix string) *Env {
	return &Env{Prefix: prefix}
}

func (p *Env) Load() (map[string]any, error) {
	out := make(map[string]any)
	prefix := p.Prefix + "_"
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		key = strings.TrimPrefix(key, prefix)
		key = strings.ToLower(strings.ReplaceAll(key, "_", "."))
		out[key] = value
	}
	return out, nil
}
