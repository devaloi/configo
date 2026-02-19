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

	server, ok := m["server"].(map[string]any)
	if !ok {
		t.Fatal("missing server key")
	}
	if server["host"] != "localhost" {
		t.Errorf("server.host = %v, want localhost", server["host"])
	}
	if m["debug"] != false {
		t.Errorf("debug = %v, want false", m["debug"])
	}
}
