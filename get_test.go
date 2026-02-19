package configo

import (
	"testing"
	"time"
)

func newTestConfig(data map[string]any) *Config {
	c := &Config{data: Flatten(data)}
	return c
}

func TestGetString(t *testing.T) {
	cfg := newTestConfig(map[string]any{"host": "localhost"})
	got, err := Get[string](cfg, "host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "localhost" {
		t.Errorf("got %q, want %q", got, "localhost")
	}
}

func TestGetInt(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": 8080})
	got, err := Get[int](cfg, "port")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 8080 {
		t.Errorf("got %d, want %d", got, 8080)
	}
}

func TestGetIntFromFloat(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": float64(8080)})
	got, err := Get[int](cfg, "port")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 8080 {
		t.Errorf("got %d, want %d", got, 8080)
	}
}

func TestGetIntFromString(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": "9090"})
	got, err := Get[int](cfg, "port")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 9090 {
		t.Errorf("got %d, want %d", got, 9090)
	}
}

func TestGetFloat64(t *testing.T) {
	cfg := newTestConfig(map[string]any{"rate": 0.75})
	got, err := Get[float64](cfg, "rate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0.75 {
		t.Errorf("got %f, want %f", got, 0.75)
	}
}

func TestGetBool(t *testing.T) {
	cfg := newTestConfig(map[string]any{"debug": true})
	got, err := Get[bool](cfg, "debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("got %v, want true", got)
	}
}

func TestGetBoolFromString(t *testing.T) {
	cfg := newTestConfig(map[string]any{"debug": "true"})
	got, err := Get[bool](cfg, "debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("got %v, want true", got)
	}
}

func TestGetDuration(t *testing.T) {
	cfg := newTestConfig(map[string]any{"timeout": "5s"})
	got, err := Get[time.Duration](cfg, "timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5*time.Second {
		t.Errorf("got %v, want 5s", got)
	}
}

func TestGetStringSlice(t *testing.T) {
	cfg := newTestConfig(map[string]any{"tags": []any{"web", "api"}})
	got, err := Get[[]string](cfg, "tags")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "web" || got[1] != "api" {
		t.Errorf("got %v, want [web api]", got)
	}
}

func TestGetIntSlice(t *testing.T) {
	cfg := newTestConfig(map[string]any{"delays": []any{1, 2, 3}})
	got, err := Get[[]int](cfg, "delays")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("got %v, want [1 2 3]", got)
	}
}

func TestGetMissingKey(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	_, err := Get[string](cfg, "missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if _, ok := err.(*KeyNotFoundError); !ok {
		t.Errorf("expected KeyNotFoundError, got %T", err)
	}
}

func TestGetOr(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	got := GetOr[string](cfg, "missing", "fallback")
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestGetOrExisting(t *testing.T) {
	cfg := newTestConfig(map[string]any{"key": "value"})
	got := GetOr[string](cfg, "key", "fallback")
	if got != "value" {
		t.Errorf("got %q, want %q", got, "value")
	}
}

func TestMustGet(t *testing.T) {
	cfg := newTestConfig(map[string]any{"key": "value"})
	got := MustGet[string](cfg, "key")
	if got != "value" {
		t.Errorf("got %q, want %q", got, "value")
	}
}

func TestMustGetPanics(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing key")
		}
	}()
	MustGet[string](cfg, "missing")
}

func TestGetNestedDotNotation(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	})
	got, err := Get[string](cfg, "database.host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "localhost" {
		t.Errorf("got %q, want %q", got, "localhost")
	}
}

func TestGetDeepNesting(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	})
	got, err := Get[string](cfg, "a.b.c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "deep" {
		t.Errorf("got %q, want %q", got, "deep")
	}
}
