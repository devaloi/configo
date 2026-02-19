# G08: configo — Type-Safe Config Library for Go

**Catalog ID:** G08 | **Size:** S | **Language:** Go
**Repo name:** `configo`
**One-liner:** A layered configuration library for Go supporting env vars, YAML, JSON, TOML with type-safe generics access, struct binding, validation, and hot reload.

---

## Why This Stands Out

- **Layered merging** — defaults → file → env → flags, each layer overrides the previous with clear precedence
- **Go generics** — type-safe `Get[T](key)` accessors, no casting, no `interface{}`
- **Struct binding with tags** — `config:"database.host" default:"localhost"` with automatic type coercion
- **Validation rules** — required, min/max, regex, custom validator functions — errors collected, not panicked
- **Hot reload** — watch config file for changes, notify subscribers via callback channels
- **Multi-format** — YAML, JSON, TOML, and `.env` files from a single unified API
- **Env prefix mapping** — `APP_DATABASE_HOST` automatically maps to `database.host`
- **Near-zero deps** — only YAML and TOML parsers; everything else is stdlib

---

## Architecture

```
configo/
├── config.go                  # Config struct, New(), Load(), layered merge logic
├── get.go                     # Generic accessors: Get[T], GetOr[T], MustGet[T]
├── bind.go                    # Struct binding: Bind(target), tag parsing
├── validate.go                # Validation engine: required, min, max, regex, custom
├── provider/
│   ├── provider.go            # Provider interface: Load() (map[string]any, error)
│   ├── defaults.go            # Default values provider
│   ├── yaml.go                # YAML file provider
│   ├── json.go                # JSON file provider
│   ├── toml.go                # TOML file provider
│   ├── env.go                 # Environment variable provider with prefix mapping
│   ├── dotenv.go              # .env file provider
│   └── flag.go                # Flag provider (stdlib flag integration)
├── watcher/
│   ├── watcher.go             # File watcher: fsnotify, debounce, reload trigger
│   └── watcher_test.go
├── errors.go                  # Typed errors: KeyNotFound, TypeMismatch, ValidationError
├── flatten.go                 # Flatten/unflatten nested maps to dot-notation keys
├── flatten_test.go
├── config_test.go             # Core config tests: load, merge, precedence
├── get_test.go                # Generic accessor tests
├── bind_test.go               # Struct binding tests
├── validate_test.go           # Validation tests
├── provider/
│   ├── yaml_test.go
│   ├── json_test.go
│   ├── toml_test.go
│   ├── env_test.go
│   └── dotenv_test.go
├── examples/
│   ├── basic/main.go          # Simple load + get
│   ├── struct/main.go         # Struct binding with validation
│   └── hotreload/main.go      # Hot reload with subscriber
├── testdata/
│   ├── config.yaml
│   ├── config.json
│   ├── config.toml
│   └── .env.test
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
├── .golangci.yml
├── LICENSE
└── README.md
```

---

## Layer Precedence

```
Priority (highest wins):
  4. Flags        (--database.host=...)
  3. Env vars     (APP_DATABASE_HOST=...)
  2. Config file  (config.yaml / config.json / config.toml)
  1. Defaults     (struct tags or SetDefault())

Load order:
  defaults → file → env → flags → merge → validate → ready
```

---

## API Surface

### Core API

```go
// Create and load
cfg := configo.New(
    configo.WithDefaults(map[string]any{"port": 8080}),
    configo.WithFile("config.yaml"),
    configo.WithEnvPrefix("APP"),
    configo.WithFlags(),
)
err := cfg.Load()

// Type-safe access
port, err := configo.Get[int](cfg, "server.port")
host := configo.GetOr[string](cfg, "server.host", "localhost")
debug := configo.MustGet[bool](cfg, "debug")

// Struct binding
type DatabaseConfig struct {
    Host     string `config:"database.host" default:"localhost" validate:"required"`
    Port     int    `config:"database.port" default:"5432" validate:"min=1,max=65535"`
    Name     string `config:"database.name" validate:"required"`
    Password string `config:"database.password" validate:"required"`
}
var db DatabaseConfig
err := cfg.Bind(&db)

// Validation
type Rule struct {
    Required bool
    Min, Max *float64
    Regex    string
    Custom   func(value any) error
}
err := cfg.Validate(map[string]configo.Rule{
    "server.port": {Required: true, Min: ptr(1.0), Max: ptr(65535.0)},
    "server.host": {Required: true},
    "database.dsn": {Required: true, Regex: `^postgres://`},
})

// Hot reload
cfg.OnChange(func(cfg *configo.Config) {
    log.Println("config reloaded")
})
cfg.Watch() // starts file watcher in background
```

### Provider Interface

```go
type Provider interface {
    Load() (map[string]any, error)
}
```

---

## Supported Types for Get[T]

| Type | Example |
|------|---------|
| `string` | `Get[string](cfg, "host")` |
| `int`, `int64` | `Get[int](cfg, "port")` |
| `float64` | `Get[float64](cfg, "rate")` |
| `bool` | `Get[bool](cfg, "debug")` |
| `time.Duration` | `Get[time.Duration](cfg, "timeout")` |
| `[]string` | `Get[[]string](cfg, "cors.origins")` |
| `[]int` | `Get[[]int](cfg, "retry.delays")` |

---

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| YAML | `gopkg.in/yaml.v3` |
| TOML | `github.com/BurntSushi/toml` |
| File watching | `github.com/fsnotify/fsnotify` |
| Testing | stdlib `testing` + `testdata/` fixtures |
| Linting | golangci-lint |

---

## Phased Build Plan

### Phase 1: Core & Providers

**1.1 — Project setup**
- `go mod init github.com/devaloi/configo`
- Create directory structure, Makefile, .gitignore, .golangci.yml
- Add LICENSE (MIT)

**1.2 — Flatten/unflatten utilities**
- `Flatten(map[string]any) map[string]any` — nested map to dot-notation (`database.host`)
- `Unflatten(map[string]any) map[string]any` — dot-notation back to nested
- Table-driven tests with edge cases (arrays, empty maps, deep nesting)

**1.3 — Provider interface and file providers**
- `Provider` interface: `Load() (map[string]any, error)`
- YAML provider: read file, unmarshal, flatten
- JSON provider: read file, unmarshal, flatten
- TOML provider: read file, unmarshal, flatten
- Tests with `testdata/` fixture files

**1.4 — Env and dotenv providers**
- Env provider: scan `os.Environ()`, filter by prefix, strip prefix, convert `_` to `.`, lowercase
- Dotenv provider: parse `.env` file (KEY=VALUE, comments, quotes)
- Tests: set env vars in test, verify mapping

**1.5 — Defaults and flag providers**
- Defaults provider: static `map[string]any`
- Flag provider: integrate with `flag.FlagSet`, parse registered flags
- Tests for both

### Phase 2: Config Core & Generics

**2.1 — Config struct and Load**
- `Config` struct: holds merged flat map, provider list, options
- `New(opts ...Option) *Config` — functional options pattern
- `Load()` — iterate providers in order, merge maps (later wins)
- Tests: verify layer precedence (defaults < file < env < flags)

**2.2 — Generic accessors**
- `Get[T](cfg, key) (T, error)` — type-safe retrieval with conversion
- `GetOr[T](cfg, key, fallback) T` — return fallback if missing
- `MustGet[T](cfg, key) T` — panic on missing (for startup)
- Type coercion: string→int, string→bool, string→duration, etc.
- Typed errors: `KeyNotFoundError`, `TypeMismatchError`
- Tests: every type, missing keys, type mismatches

**2.3 — Nested key access**
- Dot-notation traversal: `Get[string](cfg, "database.host")`
- Support both flattened and nested access patterns
- Tests: deep nesting, array indices

### Phase 3: Binding & Validation

**3.1 — Struct binding**
- `Bind(target any) error` — populate struct from config
- Tag parsing: `config:"key"` for path, `default:"value"` for fallback
- Support nested structs (recursive binding)
- Type coercion for struct fields
- Tests: simple struct, nested struct, defaults, missing required

**3.2 — Validation engine**
- Rule types: `required`, `min`, `max`, `regex`, `custom`
- `Validate(rules map[string]Rule) error` — validate current config
- Collect all errors (don't stop at first)
- `ValidationError` with per-field details
- Integrate with struct tags: `validate:"required,min=1,max=100"`
- Tests: each rule type, combined rules, custom validators, error collection

### Phase 4: Hot Reload & Polish

**4.1 — File watcher**
- `Watch()` — start fsnotify watcher on config file
- Debounce: configurable delay (default 500ms) to batch rapid changes
- On change: re-run `Load()`, call subscriber callbacks
- `OnChange(func(*Config))` — register subscriber
- `StopWatch()` — clean shutdown
- Tests: modify test file, verify callback fires, verify debounce

**4.2 — Examples**
- `basic/` — load YAML + env, print values
- `struct/` — bind to struct with validation
- `hotreload/` — watch file, print on change

**4.3 — README**
- Badges, install, quick start
- Layer precedence diagram
- API reference with examples for every feature
- Provider configuration
- Struct binding tag reference
- Validation rules reference
- Hot reload usage

**4.4 — Final checks**
- `go build ./...` clean
- `go test -race ./...` all pass
- `golangci-lint run` clean
- Fresh clone → build → test works

---

## Commit Plan

1. `feat: scaffold project with flatten utilities and tests`
2. `feat: add provider interface with YAML, JSON, TOML providers`
3. `feat: add env and dotenv providers with prefix mapping`
4. `feat: add defaults and flag providers`
5. `feat: add config core with layered merge and Load()`
6. `feat: add generic type-safe accessors Get[T], GetOr[T], MustGet[T]`
7. `feat: add struct binding with config and default tags`
8. `feat: add validation engine with required, min, max, regex, custom`
9. `feat: add hot reload with file watcher and subscribers`
10. `feat: add examples for basic, struct binding, and hot reload`
11. `docs: add README with API reference and usage guide`
12. `chore: final lint pass and cleanup`
