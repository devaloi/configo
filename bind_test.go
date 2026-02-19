package configo

import (
	"testing"
)

func TestBindSimpleStruct(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"database": map[string]any{
			"host": "db.example.com",
			"port": 5432,
			"name": "myapp",
		},
	})

	type DBConfig struct {
		Host string `config:"database.host"`
		Port int    `config:"database.port"`
		Name string `config:"database.name"`
	}

	var db DBConfig
	if err := cfg.Bind(&db); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if db.Host != "db.example.com" {
		t.Errorf("Host = %q, want %q", db.Host, "db.example.com")
	}
	if db.Port != 5432 {
		t.Errorf("Port = %d, want %d", db.Port, 5432)
	}
	if db.Name != "myapp" {
		t.Errorf("Name = %q, want %q", db.Name, "myapp")
	}
}

func TestBindDefaults(t *testing.T) {
	cfg := newTestConfig(map[string]any{})

	type ServerConfig struct {
		Host string `config:"server.host" default:"localhost"`
		Port int    `config:"server.port" default:"3000"`
	}

	var srv ServerConfig
	if err := cfg.Bind(&srv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if srv.Host != "localhost" {
		t.Errorf("Host = %q, want %q", srv.Host, "localhost")
	}
	if srv.Port != 3000 {
		t.Errorf("Port = %d, want %d", srv.Port, 3000)
	}
}

func TestBindOverridesDefaults(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"server": map[string]any{
			"host": "prodhost",
		},
	})

	type ServerConfig struct {
		Host string `config:"server.host" default:"localhost"`
		Port int    `config:"server.port" default:"3000"`
	}

	var srv ServerConfig
	if err := cfg.Bind(&srv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if srv.Host != "prodhost" {
		t.Errorf("Host = %q, want %q", srv.Host, "prodhost")
	}
	if srv.Port != 3000 {
		t.Errorf("Port = %d, want %d (default)", srv.Port, 3000)
	}
}

func TestBindBoolAndFloat(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"debug": true,
		"rate":  0.85,
	})

	type AppConfig struct {
		Debug bool    `config:"debug"`
		Rate  float64 `config:"rate"`
	}

	var app AppConfig
	if err := cfg.Bind(&app); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.Debug != true {
		t.Errorf("Debug = %v, want true", app.Debug)
	}
	if app.Rate != 0.85 {
		t.Errorf("Rate = %f, want 0.85", app.Rate)
	}
}

func TestBindNonPointer(t *testing.T) {
	cfg := newTestConfig(map[string]any{})
	type S struct{}
	var s S
	if err := cfg.Bind(s); err == nil {
		t.Fatal("expected error for non-pointer target")
	}
}

func TestBindNestedStruct(t *testing.T) {
	cfg := newTestConfig(map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"database": map[string]any{
			"host": "dbhost",
		},
	})

	type Server struct {
		Host string `config:"server.host"`
		Port int    `config:"server.port"`
	}
	type DB struct {
		Host string `config:"database.host"`
	}
	type AppConfig struct {
		Server Server
		DB     DB
	}

	var app AppConfig
	if err := cfg.Bind(&app); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.Server.Host != "localhost" {
		t.Errorf("Server.Host = %q, want %q", app.Server.Host, "localhost")
	}
	if app.DB.Host != "dbhost" {
		t.Errorf("DB.Host = %q, want %q", app.DB.Host, "dbhost")
	}
}
