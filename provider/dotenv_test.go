package provider

import (
	"testing"
)

func TestDotEnvProvider(t *testing.T) {
	p := NewDotEnv(testdataPath(".env.test"))
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := map[string]string{
		"APP_DATABASE_HOST": "db.example.com",
		"APP_DATABASE_PORT": "5432",
		"APP_DATABASE_NAME": "myapp",
		"APP_SERVER_HOST":   "localhost",
		"APP_SERVER_PORT":   "8080",
		"APP_DEBUG":         "true",
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

func TestDotEnvProviderMissingFile(t *testing.T) {
	p := NewDotEnv("nonexistent.env")
	_, err := p.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
