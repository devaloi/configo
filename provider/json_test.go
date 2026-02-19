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

	server, ok := m["server"].(map[string]any)
	if !ok {
		t.Fatal("missing server key")
	}
	if server["host"] != "localhost" {
		t.Errorf("server.host = %v, want localhost", server["host"])
	}
	if server["port"] != float64(8080) {
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

func TestJSONProviderMissingFile(t *testing.T) {
	p := NewJSON("nonexistent.json")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
