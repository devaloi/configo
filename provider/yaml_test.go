package provider

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataPath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata", name)
}

func TestYAMLProvider(t *testing.T) {
	p := NewYAML(testdataPath("config.yaml"))
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := map[string]any{
		"server.host":  "localhost",
		"server.port":  8080,
		"database.host": "db.example.com",
		"database.port": 5432,
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

func TestYAMLProviderMissingFile(t *testing.T) {
	p := NewYAML("nonexistent.yaml")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
