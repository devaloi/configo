package provider

import (
	"os"
	"testing"
)

func TestEnvProvider(t *testing.T) {
	os.Setenv("MYAPP_DATABASE_HOST", "envhost")
	os.Setenv("MYAPP_DATABASE_PORT", "3306")
	os.Setenv("MYAPP_DEBUG", "true")
	defer func() {
		os.Unsetenv("MYAPP_DATABASE_HOST")
		os.Unsetenv("MYAPP_DATABASE_PORT")
		os.Unsetenv("MYAPP_DEBUG")
	}()

	p := NewEnv("MYAPP")
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := map[string]string{
		"database.host": "envhost",
		"database.port": "3306",
		"debug":         "true",
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

func TestEnvProviderNoMatch(t *testing.T) {
	os.Setenv("OTHER_KEY", "value")
	defer os.Unsetenv("OTHER_KEY")

	p := NewEnv("MYAPP")
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["key"]; ok {
		t.Error("should not include non-prefixed keys")
	}
}
