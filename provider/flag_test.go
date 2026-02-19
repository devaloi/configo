package provider

import (
	"flag"
	"testing"
)

func TestFlagProvider(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("server.host", "", "server host")
	fs.Int("server.port", 0, "server port")

	err := fs.Parse([]string{"--server.host=flaghost", "--server.port=9090"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	p := NewFlag(fs)
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := m["server.host"]; got != "flaghost" {
		t.Errorf("server.host = %v, want flaghost", got)
	}
	if got := m["server.port"]; got != "9090" {
		t.Errorf("server.port = %v, want 9090", got)
	}
}

func TestFlagProviderUnsetFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("server.host", "default", "server host")
	_ = fs.Parse([]string{})

	p := NewFlag(fs)
	m, err := p.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := m["server.host"]; ok {
		t.Error("unset flags should not be included")
	}
}
