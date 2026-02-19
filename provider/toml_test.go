package provider

import (
	"testing"
)

func TestTOMLProvider(t *testing.T) {
	p := NewTOML(testdataPath("config.toml"))
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
	if server["port"] != int64(8080) {
		t.Errorf("server.port = %v (%T), want 8080", server["port"], server["port"])
	}

	db, ok := m["database"].(map[string]any)
	if !ok {
		t.Fatal("missing database key")
	}
	if db["host"] != "db.example.com" {
		t.Errorf("database.host = %v, want db.example.com", db["host"])
	}
}

func TestTOMLProviderMissingFile(t *testing.T) {
	p := NewTOML("nonexistent.toml")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
