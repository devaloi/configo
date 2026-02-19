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

	server, ok := m["server"].(map[string]any)
	if !ok {
		t.Fatal("missing server key")
	}
	if server["host"] != "localhost" {
		t.Errorf("server.host = %v, want localhost", server["host"])
	}
	if server["port"] != 8080 {
		t.Errorf("server.port = %v, want 8080", server["port"])
	}

	db, ok := m["database"].(map[string]any)
	if !ok {
		t.Fatal("missing database key")
	}
	if db["host"] != "db.example.com" {
		t.Errorf("database.host = %v, want db.example.com", db["host"])
	}
}

func TestYAMLProviderMissingFile(t *testing.T) {
	p := NewYAML("nonexistent.yaml")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
