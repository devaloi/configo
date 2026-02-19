package provider

import (
	"flag"
	"fmt"
)

// Flag loads configuration from a flag.FlagSet.
type Flag struct {
	FlagSet *flag.FlagSet
}

func NewFlag(fs *flag.FlagSet) *Flag {
	return &Flag{FlagSet: fs}
}

func (p *Flag) Load() (map[string]any, error) {
	out := make(map[string]any)
	p.FlagSet.Visit(func(f *flag.Flag) {
		out[f.Name] = f.Value.String()
	})
	return out, nil
}

// RegisterFlags defines common config flags on a FlagSet for convenience.
func RegisterFlags(fs *flag.FlagSet, keys ...string) {
	for _, key := range keys {
		fs.String(key, "", fmt.Sprintf("config value for %s", key))
	}
}
