package configo

import (
	"flag"
	"os"
	"testing"
)

func TestConfigLoadYAML(t *testing.T) {
	cfg := New(
		WithFile("testdata/config.yaml"),
	)
	if err := cfg.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := cfg.Data()
	if got := data["server.host"]; got != "localhost" {
		t.Errorf("server.host = %v, want localhost", got)
	}
}

func TestConfigLoadJSON(t *testing.T) {
	cfg := New(
		WithFile("testdata/config.json"),
	)
	if err := cfg.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := cfg.Data()
	if got := data["server.host"]; got != "localhost" {
		t.Errorf("server.host = %v, want localhost", got)
	}
}

func TestConfigLoadTOML(t *testing.T) {
	cfg := New(
		WithFile("testdata/config.toml"),
	)
	if err := cfg.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := cfg.Data()
	if got := data["server.host"]; got != "localhost" {
		t.Errorf("server.host = %v, want localhost", got)
	}
}

func TestConfigLayerPrecedence(t *testing.T) {
	os.Setenv("TEST_SERVER_HOST", "envhost")
	defer os.Unsetenv("TEST_SERVER_HOST")

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("server.host", "", "host")
	_ = fs.Parse([]string{"--server.host=flaghost"})

	cfg := New(
		WithDefaults(map[string]any{"server.host": "defaulthost", "server.port": 3000}),
		WithFile("testdata/config.yaml"),
		WithEnvPrefix("TEST"),
		WithFlags(fs),
	)
	if err := cfg.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := cfg.Data()
	// Flags have highest priority
	if got := data["server.host"]; got != "flaghost" {
		t.Errorf("server.host = %v, want flaghost (flags override all)", got)
	}
	// File overrides default
	if got := data["server.port"]; got != 8080 {
		t.Errorf("server.port = %v, want 8080 (file overrides default)", got)
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg := New(
		WithDefaults(map[string]any{
			"app.name": "testapp",
			"app.port": 9090,
		}),
	)
	if err := cfg.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := cfg.Data()
	if got := data["app.name"]; got != "testapp" {
		t.Errorf("app.name = %v, want testapp", got)
	}
}
