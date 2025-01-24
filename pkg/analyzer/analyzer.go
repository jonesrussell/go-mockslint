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
	moduleFileName = "module.go"
	internalDir    = "internal"
	moduleDir      = "module"
)

// Config holds the analyzer configuration
type Config struct {
	ModulePaths  []string `yaml:"modulePaths"`
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
	Name: "fxlint",
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
	modulePaths := stringSliceFlag{
		"internal/*/module.go",
		"pkg/*/module.go",
	}

	Analyzer.Flags.Var(&modulePaths, "module-paths", "Allowed module file paths (glob patterns)")
	Analyzer.Flags.BoolVar(&config.StrictNaming, "strict-naming", true, "Enforce strict module naming")
	config.ModulePaths = modulePaths
}

func init() {
	config = Config{
		ModulePaths: []string{
			"internal/*/module.go",
			"pkg/*/module.go",
		},
		StrictNaming: true,
	}

	setupFlags()
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// We only care about fx.Module calls
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		checkModuleCall(pass, n.(*ast.CallExpr))
	})

	return nil, nil
}

func checkModuleCall(pass *analysis.Pass, call *ast.CallExpr) {
	// Check if it's an fx.Module call
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, isIdent := sel.X.(*ast.Ident)
	if !isIdent || ident.Name != "fx" || sel.Sel.Name != "Module" {
		return
	}

	checkModuleLocation(pass, call)
}

func checkModuleLocation(pass *analysis.Pass, call *ast.CallExpr) {
	// Get file info
	pos := call.Pos()
	file := pass.Fset.File(pos)
	filename := filepath.Base(file.Name())
	dir := filepath.Dir(file.Name())
	parts := strings.Split(dir, string(filepath.Separator))

	// Check if file is module.go
	if filename != moduleFileName {
		// Only report on init functions
		if fn, isInit := findParentInit(call); isInit {
			pass.Reportf(fn.Pos(), "fx.Module can only be used in module.go files")
		}

		return
	}

	// Check if file is in internal/ or internal/module/
	if isInvalidLocation(parts) {
		pass.Reportf(call.Pos(), "module.go files should not be directly in internal/ or internal/module/ directories")

		return
	}

	checkModuleName(pass, call, dir)
}

func checkModuleName(pass *analysis.Pass, call *ast.CallExpr, dir string) {
	if len(call.Args) == 0 {
		return
	}

	lit, isLit := call.Args[0].(*ast.BasicLit)
	if !isLit {
		return
	}

	moduleName := strings.Trim(lit.Value, `"`)
	dirName := filepath.Base(dir)

	// Check against package name
	if file := findEnclosingFile(pass, call); file != nil && file.Name != nil {
		pkgName := file.Name.Name
		if moduleName != pkgName {
			pass.Reportf(lit.Pos(), "module name %q should match package name %q", moduleName, pkgName)

			return
		}
	}

	// Check against directory name
	if moduleName != dirName {
		pass.Reportf(lit.Pos(), "module name %q should match directory name %q", moduleName, dirName)
	}
}

func isInvalidLocation(parts []string) bool {
	for i, part := range parts {
		if part == internalDir {
			return i == len(parts)-1 || parts[i+1] == moduleDir
		}
	}

	return false
}

func findEnclosingFile(pass *analysis.Pass, node ast.Node) *ast.File {
	for _, file := range pass.Files {
		if file.Pos() <= node.Pos() && node.Pos() <= file.End() {
			return file
		}
	}

	return nil
}

func findParentInit(node ast.Node) (*ast.FuncDecl, bool) {
	var initFunc *ast.FuncDecl

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "init" {
			initFunc = fn

			return false
		}

		return true
	})

	return initFunc, initFunc != nil
}
