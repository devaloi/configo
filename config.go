package configo

import (
	"flag"
	"sync"

	"github.com/devaloi/configo/provider"
)

// Config holds merged configuration data and providers.
type Config struct {
	mu        sync.RWMutex
	data      map[string]any
	providers []provider.Provider
	filePath  string
	onChange  []func(*Config)
}

// Option configures a Config instance.
type Option func(*Config)

// New creates a new Config with the given options.
func New(opts ...Option) *Config {
	c := &Config{
		data: make(map[string]any),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithDefaults adds a defaults provider (lowest priority).
func WithDefaults(defaults map[string]any) Option {
	return func(c *Config) {
		c.providers = append(c.providers, provider.NewDefaults(defaults))
	}
}

// WithFile adds a file provider based on extension (.yaml/.yml, .json, .toml).
func WithFile(path string) Option {
	return func(c *Config) {
		c.filePath = path
		var p provider.Provider
		switch {
		case hasExt(path, ".yaml", ".yml"):
			p = provider.NewYAML(path)
		case hasExt(path, ".json"):
			p = provider.NewJSON(path)
		case hasExt(path, ".toml"):
			p = provider.NewTOML(path)
		}
		if p != nil {
			c.providers = append(c.providers, p)
		}
	}
}

// WithEnvPrefix adds an environment variable provider.
func WithEnvPrefix(prefix string) Option {
	return func(c *Config) {
		c.providers = append(c.providers, provider.NewEnv(prefix))
	}
}

// WithDotEnv adds a .env file provider.
func WithDotEnv(path string) Option {
	return func(c *Config) {
		c.providers = append(c.providers, provider.NewDotEnv(path))
	}
}

// WithFlags adds a flag provider using the given FlagSet.
func WithFlags(fs *flag.FlagSet) Option {
	return func(c *Config) {
		c.providers = append(c.providers, provider.NewFlag(fs))
	}
}

// WithProvider adds a custom provider.
func WithProvider(p provider.Provider) Option {
	return func(c *Config) {
		c.providers = append(c.providers, p)
	}
}

// Load iterates all providers in order and merges their data.
// Later providers override earlier ones.
func (c *Config) Load() error {
	merged := make(map[string]any)
	for _, p := range c.providers {
		m, err := p.Load()
		if err != nil {
			return err
		}
		flat := Flatten(m)
		for k, v := range flat {
			merged[k] = v
		}
	}
	c.mu.Lock()
	c.data = merged
	c.mu.Unlock()
	return nil
}

// Data returns a copy of the current configuration data.
func (c *Config) Data() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]any, len(c.data))
	for k, v := range c.data {
		out[k] = v
	}
	return out
}

// OnChange registers a callback that fires when config is reloaded.
func (c *Config) OnChange(fn func(*Config)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onChange = append(c.onChange, fn)
}

func hasExt(path string, exts ...string) bool {
	for _, ext := range exts {
		if len(path) > len(ext) && path[len(path)-len(ext):] == ext {
			return true
		}
	}
	return false
}
