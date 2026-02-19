# configo

[![Go](https://img.shields.io/badge/Go-1.26-blue.svg)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A layered configuration library for Go supporting env vars, YAML, JSON, TOML with type-safe generics access, struct binding, validation, and hot reload.

## Features

- **Layered merging** — defaults → file → env → flags, each layer overrides the previous
- **Go generics** — type-safe `Get[T](key)` accessors, no casting
- **Struct binding** — `config:"key"` and `default:"value"` tags with automatic type coercion
- **Validation** — required, min/max, regex, custom validators; errors collected, not panicked
- **Hot reload** — watch config files for changes, notify subscribers
- **Multi-format** — YAML, JSON, TOML, and `.env` files from a single API
- **Env prefix mapping** — `APP_DATABASE_HOST` → `database.host`

## Install

```bash
go get github.com/devaloi/configo
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/devaloi/configo"
)

func main() {
    cfg := configo.New(
        configo.WithDefaults(map[string]any{"server.port": 3000}),
        configo.WithFile("config.yaml"),
        configo.WithEnvPrefix("APP"),
    )
    if err := cfg.Load(); err != nil {
        log.Fatal(err)
    }

    host := configo.GetOr[string](cfg, "server.host", "localhost")
    port := configo.GetOr[int](cfg, "server.port", 3000)
    fmt.Printf("Server: %s:%d\n", host, port)
}
```

## Layer Precedence

```
Priority (highest wins):
  4. Flags        (--database.host=...)
  3. Env vars     (APP_DATABASE_HOST=...)
  2. Config file  (config.yaml / config.json / config.toml)
  1. Defaults     (WithDefaults or struct tags)
```

## API Reference

### Creating a Config

```go
cfg := configo.New(
    configo.WithDefaults(map[string]any{...}),
    configo.WithFile("config.yaml"),
    configo.WithEnvPrefix("APP"),
    configo.WithDotEnv(".env"),
    configo.WithFlags(flagSet),
    configo.WithProvider(customProvider),
)
err := cfg.Load()
```

### Type-Safe Accessors

```go
// Get with error handling
port, err := configo.Get[int](cfg, "server.port")

// Get with fallback
host := configo.GetOr[string](cfg, "server.host", "localhost")

// Get or panic (for startup)
debug := configo.MustGet[bool](cfg, "debug")
```

#### Supported Types

| Type | Example |
|------|---------|
| `string` | `Get[string](cfg, "host")` |
| `int`, `int64` | `Get[int](cfg, "port")` |
| `float64` | `Get[float64](cfg, "rate")` |
| `bool` | `Get[bool](cfg, "debug")` |
| `time.Duration` | `Get[time.Duration](cfg, "timeout")` |
| `[]string` | `Get[[]string](cfg, "cors.origins")` |
| `[]int` | `Get[[]int](cfg, "retry.delays")` |

### Struct Binding

```go
type DatabaseConfig struct {
    Host     string `config:"database.host" default:"localhost" validate:"required"`
    Port     int    `config:"database.port" default:"5432"      validate:"min=1,max=65535"`
    Name     string `config:"database.name"                     validate:"required"`
    Password string `config:"database.password"`
}

var db DatabaseConfig
err := cfg.Bind(&db)
```

#### Tag Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `config` | Config key path | `config:"database.host"` |
| `default` | Fallback value | `default:"localhost"` |
| `validate` | Validation rules | `validate:"required,min=1"` |

### Validation

```go
// Map-based rules
err := cfg.Validate(map[string]configo.Rule{
    "server.port": {Required: true, Min: ptr(1.0), Max: ptr(65535.0)},
    "server.host": {Required: true},
    "database.dsn": {Regex: `^postgres://`},
    "app.name": {Custom: func(v any) error {
        if v == "" {
            return fmt.Errorf("must not be empty")
        }
        return nil
    }},
})

// Struct-based validation (uses validate tags)
err := cfg.ValidateStruct(DatabaseConfig{})
```

#### Validation Rules

| Rule | Description | Example |
|------|-------------|---------|
| `required` | Value must exist | `validate:"required"` |
| `min=N` | Minimum numeric value | `validate:"min=1"` |
| `max=N` | Maximum numeric value | `validate:"max=65535"` |
| `regex=PATTERN` | Must match regex | `validate:"regex=^https://"` |
| Custom | Custom function | `Rule{Custom: fn}` |

Validation collects all errors into a `ValidationError` — it does not stop at the first failure.

### Hot Reload

```go
cfg.OnChange(func(c *configo.Config) {
    log.Println("config reloaded")
})
err := cfg.Watch()    // starts file watcher in background
defer cfg.StopWatch() // clean shutdown
```

### Providers

All providers implement the `Provider` interface:

```go
type Provider interface {
    Load() (map[string]any, error)
}
```

| Provider | Description |
|----------|-------------|
| `provider.NewDefaults(map)` | Static default values |
| `provider.NewYAML(path)` | YAML file |
| `provider.NewJSON(path)` | JSON file |
| `provider.NewTOML(path)` | TOML file |
| `provider.NewEnv(prefix)` | Environment variables with prefix |
| `provider.NewDotEnv(path)` | `.env` file |
| `provider.NewFlag(flagSet)` | stdlib `flag.FlagSet` |

### Env Variable Mapping

With `WithEnvPrefix("APP")`:

```
APP_DATABASE_HOST  →  database.host
APP_SERVER_PORT    →  server.port
APP_DEBUG          →  debug
```

## Error Types

| Error | Description |
|-------|-------------|
| `KeyNotFoundError` | Requested key does not exist |
| `TypeMismatchError` | Value cannot be converted to requested type |
| `ValidationError` | One or more validation rules failed (contains `[]FieldError`) |

## Examples

See the [`examples/`](examples/) directory:

- [`basic/`](examples/basic/) — Load config and read values
- [`struct/`](examples/struct/) — Struct binding with validation
- [`hotreload/`](examples/hotreload/) — Watch file for changes

## License

[MIT](LICENSE)
