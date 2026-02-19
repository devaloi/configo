# Build configo — Type-Safe Config Library for Go

You are building a **portfolio project** for a Senior AI Engineer's public GitHub. It must be impressive, clean, and production-grade. Read these docs before writing any code:

1. **`G08-go-config-library.md`** — Complete project spec: architecture, phases, provider system, generic accessors, validation, hot reload, commit plan. This is your primary blueprint. Follow it phase by phase.
2. **`github-portfolio.md`** — Portfolio goals and Definition of Done (Level 1 + Level 2). Understand the quality bar.
3. **`github-portfolio-checklist.md`** — Pre-publish checklist. Every item must pass before you're done.

---

## Instructions

### Read first, build second
Read all three docs completely before writing a single line of code. Understand the layered config architecture, the provider interface, the generic accessor system, struct binding with tags, the validation engine, and the hot reload watcher.

### Follow the phases in order
The project spec has 4 phases. Do them in order:
1. **Core & Providers** — project setup, flatten/unflatten utilities, provider interface, YAML/JSON/TOML/env/dotenv/flag providers with tests
2. **Config Core & Generics** — Config struct with layered merge, generic Get[T]/GetOr[T]/MustGet[T] accessors, type coercion, typed errors
3. **Binding & Validation** — struct binding with `config:` and `default:` tags, validation engine with required/min/max/regex/custom rules, error collection
4. **Hot Reload & Polish** — fsnotify file watcher with debounce, subscriber callbacks, examples, comprehensive README

### Commit frequently
Follow the commit plan in the spec. Use **conventional commits**. Each commit should be a logical unit.

### Quality non-negotiables
- **Go generics for type safety.** `Get[T]` must use generics — no `interface{}` returns, no type assertion by the caller. The library handles all type coercion internally.
- **Layered precedence is the core feature.** defaults → file → env → flags must merge correctly. A value set in env must override the same key from a file. Tests must prove this.
- **Provider interface.** Every source (YAML, JSON, TOML, env, flags) implements the same `Provider` interface. Adding a new source is one file.
- **Struct binding with tags.** `config:"database.host" default:"localhost" validate:"required"` must all work together. Nested struct binding must work recursively.
- **Validation collects all errors.** Don't stop at the first validation failure. Return a `ValidationError` listing every field that failed and why.
- **Hot reload is real.** fsnotify watches the file, debounces rapid changes, re-runs Load(), and calls all registered subscribers. Tests must verify this.
- **Near-zero dependencies.** Only `gopkg.in/yaml.v3`, `github.com/BurntSushi/toml`, and `github.com/fsnotify/fsnotify`. Everything else is stdlib.
- **Lint clean.** `golangci-lint run` and `go vet` must pass with zero warnings.

### What NOT to do
- Don't use Viper or any config library. This IS the config library, built from scratch.
- Don't return `interface{}` from accessors. Use Go generics for type-safe returns.
- Don't panic on missing keys. Return typed errors (`KeyNotFoundError`, `TypeMismatchError`). Only `MustGet` panics, and that's intentional.
- Don't skip struct binding tests. Nested structs, default values, and validation tags must all be tested.
- Don't leave `// TODO` or `// FIXME` comments anywhere.
- Don't commit testdata that contains real secrets or credentials.

---

## GitHub Username

The GitHub username is **devaloi**. For Go module paths, use `github.com/devaloi/configo`. All internal imports must use this module path.

## Start

Read the three docs. Then begin Phase 1 from `G08-go-config-library.md`.
