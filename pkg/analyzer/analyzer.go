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

var Analyzer = &analysis.Analyzer{
	Name: "fxlint",
	Doc:  "Enforces domain-driven module organization patterns in Go projects using uber/fx",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.File)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			checkModuleFile(pass, node)
		case *ast.CallExpr:
			checkModuleCall(pass, node)
		}
	})

	return nil, nil
}

func checkModuleFile(pass *analysis.Pass, file *ast.File) {
	filename := filepath.Base(pass.Fset.File(file.Pos()).Name())
	if filename != moduleFileName {
		return
	}

	// Check if file is directly in internal directory
	dir := filepath.Dir(pass.Fset.File(file.Pos()).Name())
	parts := strings.Split(dir, string(filepath.Separator))
	for i, part := range parts {
		if part == internalDir {
			if i == len(parts)-1 || parts[i+1] == moduleDir {
				pass.Reportf(file.Pos(), "module.go files should not be directly in internal/ or internal/module/ directories")
			}
		}
	}

	// Check if module name matches package name
	if file.Name != nil {
		pkgName := file.Name.Name
		checkModuleNameMatchesPackage(pass, file, pkgName)
	}
}

func checkModuleCall(pass *analysis.Pass, call *ast.CallExpr) {
	// Check for fx.Module calls
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if x, ok := sel.X.(*ast.Ident); ok && x.Name == "fx" && sel.Sel.Name == "Module" {
			filename := filepath.Base(pass.Fset.File(call.Pos()).Name())
			if filename != moduleFileName {
				pass.Reportf(call.Pos(), "fx.Module can only be used in module.go files")
			}

			// Check module name argument
			if len(call.Args) > 0 {
				if lit, ok := call.Args[0].(*ast.BasicLit); ok {
					moduleName := strings.Trim(lit.Value, `"`)
					dir := filepath.Base(filepath.Dir(pass.Fset.File(call.Pos()).Name()))
					if moduleName != dir {
						pass.Reportf(lit.Pos(), "module name %q should match directory name %q", moduleName, dir)
					}
				}
			}
		}
	}
}

func checkModuleNameMatchesPackage(pass *analysis.Pass, file *ast.File, pkgName string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if x, ok := sel.X.(*ast.Ident); ok && x.Name == "fx" && sel.Sel.Name == "Module" {
					if len(call.Args) > 0 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							moduleName := strings.Trim(lit.Value, `"`)
							if moduleName != pkgName {
								pass.Reportf(lit.Pos(), "module name %q should match package name %q", moduleName, pkgName)
							}
						}
					}
				}
			}
		}
		return true
	})
}
