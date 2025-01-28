# go-mockslint

A custom golangci-lint analyzer that enforces mock organization patterns in Go projects. The linter ensures that:

1. Mocks are only created in `test/mocks/` directory
2. No mocks are allowed in `internal/` directories
3. Proper mock naming conventions are followed

## Installation

```bash
go install github.com/jonesrussell/go-mockslint/cmd/go-mockslint@latest
```

## Features

The linter enforces the following rules:

1. Mock files must be in the `test/mocks/` directory
2. No mock files allowed in `internal/` directories
3. Mock files should follow naming convention: `mock_*.go`
4. Each mock should be in its own file
5. Mock interfaces should be properly documented

## Usage

### As a standalone tool

```bash
# Run on a package
go-mockslint ./...

# Run with custom configuration
go-mockslint -mock-paths="test/mocks/*_test.go" -strict-naming=true ./...
```

### With golangci-lint

Add to your `.golangci.yml`:

```yaml
linters:
  enable:
    - mockslint

linters-settings:
  mockslint:
    # Configure allowed mock file locations
    mockPaths:
      - test/mocks/*.go
    # Enforce strict mock naming
    strictNaming: true
```

Then run:

```bash
golangci-lint run
```

## Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `mockPaths` | Glob patterns for allowed mock file locations | `["test/mocks/*.go"]` |
| `strictNaming` | Enforce that mock files follow naming convention | `true` |

## Examples

### Valid Mock Organization

```go
// test/mocks/mock_authenticator.go
package mocks

import "github.com/stretchr/testify/mock"

type MockAuthenticator struct {
    mock.Mock
}

func (m *MockAuthenticator) Authenticate(token string) bool {
    args := m.Called(token)
    return args.Bool(0)
}
```

### Common Violations

```go
// BAD: Mock in wrong location
// internal/auth/mock_auth.go
package auth

type MockAuth struct {}

// BAD: Mock not in test/mocks
// pkg/auth/mocks/authenticator.go
package mocks

type MockAuthenticator struct {}

// BAD: Wrong naming convention
// test/mocks/auth_mock.go
package mocks

type AuthMock struct {}
```

## Development

### Prerequisites

- Go 1.23 or later
- [Task](https://taskfile.dev) for running development commands
- [golangci-lint](https://golangci-lint.run/) for linting

### Commands

```bash
# Build the binary
task build

# Run tests
task test

# Run linter
task lint

# Run all CI checks
task ci

# Install globally
task install

# Clean build artifacts
task clean
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
name: Tests
on:
  push:
    branches:
      - main
      - master
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: arduino/setup-task@v2
        with:
          version: 3.x
      - run: task ci
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the [LICENSE](LICENSE) file for details 