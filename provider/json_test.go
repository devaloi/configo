package provider

import (
	"testing"
)

func TestJSONProvider(t *testing.T) {
	p := NewJSON(testdataPath("config.json"))
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := map[string]any{
		"server.host":   "localhost",
		"server.port":   float64(8080),
		"database.host": "db.example.com",
		"database.port": float64(5432),
		"database.name": "myapp",
		"server.debug":  true,
	}
	for key, want := range tests {
		got, ok := m[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if got != want {
			t.Errorf("key %q = %v (%T), want %v (%T)", key, got, got, want, want)
		}
	}
}

func TestJSONProviderMissingFile(t *testing.T) {
	p := NewJSON("nonexistent.json")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
