package analyzer

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	mocksDir    = "test/mocks"
	internalDir = "internal"
)

// Config holds the analyzer configuration
type Config struct {
	MockPaths    []string `yaml:"mockPaths"`
	StrictNaming bool     `yaml:"strictNaming"`
}

// stringSliceFlag implements the flag.Value interface.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// Analyzer is the fxlint analyzer.
var Analyzer = &analysis.Analyzer{
	Name: "mockslint",
	Doc:  "Enforces domain-driven module organization patterns in Go projects using uber/fx",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// config holds the analyzer configuration.
var config Config

// setupFlags initializes the analyzer flags.
func setupFlags() {
	mockPaths := stringSliceFlag{
		"test/mocks/*.go",
	}

	Analyzer.Flags.Var(&mockPaths, "mock-paths", "Allowed mock file paths (glob patterns)")
	Analyzer.Flags.BoolVar(&config.StrictNaming, "strict-naming", true, "Enforce strict mock naming")
	config.MockPaths = mockPaths
}

func init() {
	config = Config{
		MockPaths: []string{
			"test/mocks/*.go",
		},
		StrictNaming: true,
	}

	setupFlags()
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec := n.(*ast.TypeSpec)
		if strings.HasPrefix(typeSpec.Name.Name, "Mock") {
			checkMockLocation(pass, typeSpec)
		}
	})

	return nil, nil
}

func checkMockLocation(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	pos := typeSpec.Pos()
	file := pass.Fset.File(pos)
	dir := filepath.Dir(file.Name())
	parts := strings.Split(dir, string(filepath.Separator))

	// Check if mock is in internal/
	if containsDir(parts, internalDir) {
		pass.Reportf(pos, "mock types are not allowed in internal/ directories")
		return
	}

	// Check if mock is in test/mocks/
	if !strings.HasPrefix(dir, mocksDir) {
		pass.Reportf(pos, "mock types must be defined in test/mocks/ directory")
		return
	}
}

func containsDir(parts []string, dir string) bool {
	for _, part := range parts {
		if part == dir {
			return true
		}
	}
	return false
}
