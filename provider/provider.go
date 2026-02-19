package provider

// Provider loads configuration from a source and returns a flat map.
type Provider interface {
	Load() (map[string]any, error)
}
