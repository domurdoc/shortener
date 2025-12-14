package exitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  `check for os.Exit calls within main package main func`,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if !isMainPackageFile(file) {
			continue
		}
		mainFunc := getMainFuncDecl(file)
		if mainFunc == nil {
			continue
		}
		ast.Inspect(mainFunc, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if selIdent, ok := sel.X.(*ast.Ident); ok {
						if sel.Sel.Name == "Exit" && selIdent.Name == "os" {
							pass.Reportf(n.Pos(), "os.Exit call within main package main func")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func isMainPackageFile(f *ast.File) bool {
	return f.Name.Name == "main"
}

func getMainFuncDecl(f *ast.File) *ast.FuncDecl {
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.FuncDecl); ok {
			if d.Name.Name == "main" {
				return d
			}
		}
	}
	return nil
}
