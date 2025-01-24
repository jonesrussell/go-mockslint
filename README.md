# go-fxlint

A custom golangci-lint analyzer that enforces domain-driven module organization patterns in Go projects using [uber/fx](https://github.com/uber-go/fx).

## Installation

```bash
go install github.com/jonesrussell/go-fxlint/cmd/go-fxlint@latest
```

## Features

The linter enforces the following rules:

1. Module files must be named `module.go`
2. Module files must be in their respective domain packages (not in `internal/module/`)
3. Module files cannot be directly in the `internal/` directory
4. The fx.Module name should match its package name (e.g., package auth should have `fx.Module("auth", ...)`)
5. `fx.Module` can only be used in `module.go` files
6. Each domain package should have its own `module.go`

## Usage

### As a standalone tool

```bash
# Run on a package
go-fxlint ./...

# Run with custom configuration
go-fxlint -module-paths="internal/*/module.go,pkg/*/module.go" -strict-naming=true ./...
```

### With golangci-lint

Add to your `.golangci.yml`:

```yaml
linters:
  enable:
    - fxlint

linters-settings:
  fxlint:
    # Optional: Configure allowed module locations
    module-paths:
      - internal/*/module.go
      - pkg/*/module.go
    # Optional: Enforce strict module naming
    strict-naming: true
```

Then run:

```bash
golangci-lint run
```

## Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `module-paths` | Glob patterns for allowed module file locations | `["internal/*/module.go", "pkg/*/module.go"]` |
| `strict-naming` | Enforce that module names match their package names | `true` |

## Examples

### Valid Module Organization

```go
// internal/auth/module.go
package auth

import "go.uber.org/fx"

var Module = fx.Module("auth",
    fx.Provide(
        NewAuthenticator,
    ),
)
```

### Common Violations

```go
// BAD: Module in wrong location
// internal/module/auth.go
package auth

var Module = fx.Module("auth", ...)

// BAD: Module directly in internal
// internal/module.go
package internal

var Module = fx.Module("internal", ...)

// BAD: Wrong module name
// internal/auth/module.go
package auth

var Module = fx.Module("user", ...) // Should be "auth"

// BAD: fx.Module in non-module file
// internal/auth/service.go
package auth

func init() {
    fx.Module("auth", ...) // Should be in module.go
}
```

## Integration

### VSCode

Add to your `.vscode/settings.json`:

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ]
}
```

### GitHub Actions

```yaml
name: Lint
on: [push, pull_request]
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the [LICENSE](LICENSE) file for details 