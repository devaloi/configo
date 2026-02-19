package provider

import (
	"testing"
)

func TestDefaultsProvider(t *testing.T) {
	p := NewDefaults(map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"debug": false,
	})
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := map[string]any{
		"server.host": "localhost",
		"server.port": 8080,
		"debug":       false,
	}
	for key, want := range tests {
		got, ok := m[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if got != want {
			t.Errorf("key %q = %v, want %v", key, got, want)
		}
	}
}
