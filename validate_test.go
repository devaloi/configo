package configo

import (
	"errors"
	"fmt"
	"testing"
)

func ptr(f float64) *float64 { return &f }

func TestValidateRequired(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	err := cfg.Validate(map[string]Rule{
		"server.host": {Required: true},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(ve.Errors))
	}
}

func TestValidateRequiredPresent(t *testing.T) {
	cfg := newTestConfig(map[string]any{"server.host": "localhost"})
	err := cfg.Validate(map[string]Rule{
		"server.host": {Required: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMin(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": 0})
	err := cfg.Validate(map[string]Rule{
		"port": {Min: ptr(1)},
	})
	if err == nil {
		t.Fatal("expected validation error for min")
	}
}

func TestValidateMax(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": 70000})
	err := cfg.Validate(map[string]Rule{
		"port": {Max: ptr(65535)},
	})
	if err == nil {
		t.Fatal("expected validation error for max")
	}
}

func TestValidateMinMax(t *testing.T) {
	cfg := newTestConfig(map[string]any{"port": 8080})
	err := cfg.Validate(map[string]Rule{
		"port": {Min: ptr(1), Max: ptr(65535)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRegex(t *testing.T) {
	cfg := newTestConfig(map[string]any{"dsn": "mysql://localhost"})
	err := cfg.Validate(map[string]Rule{
		"dsn": {Regex: `^postgres://`},
	})
	if err == nil {
		t.Fatal("expected validation error for regex")
	}
}

func TestValidateRegexMatch(t *testing.T) {
	cfg := newTestConfig(map[string]any{"dsn": "postgres://localhost"})
	err := cfg.Validate(map[string]Rule{
		"dsn": {Regex: `^postgres://`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCustom(t *testing.T) {
	cfg := newTestConfig(map[string]any{"name": ""})
	err := cfg.Validate(map[string]Rule{
		"name": {Custom: func(v any) error {
			if v == "" {
				return fmt.Errorf("must not be empty")
			}
			return nil
		}},
	})
	if err == nil {
		t.Fatal("expected validation error for custom")
	}
}

func TestValidateMultipleErrors(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	err := cfg.Validate(map[string]Rule{
		"host": {Required: true},
		"port": {Required: true},
	})
	if err == nil {
		t.Fatal("expected validation errors")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(ve.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(ve.Errors))
	}
}

func TestValidateStruct(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	})

	type DBConfig struct {
		Host string `config:"database.host" validate:"required"`
		Port int    `config:"database.port" validate:"required,min=1,max=65535"`
		Name string `config:"database.name" validate:"required"`
	}

	err := cfg.ValidateStruct(DBConfig{})
	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	found := false
	for _, fe := range ve.Errors {
		if fe.Field == "database.name" {
			found = true
		}
	}
	if !found {
		t.Error("expected error for database.name")
	}
}
