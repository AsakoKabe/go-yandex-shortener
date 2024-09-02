package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const doc = "check for use of os.Exit in main"

// Analyzer запрещающий использовать прямой вызов os.Exit в функции main пакета main
var Analyzer = &analysis.Analyzer{
	Name: "no-os-exit",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(
			file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.FuncDecl:
					if x.Name.Name != "main" || pass.Pkg.Name() != "main" {
						return true
					}

					ast.Inspect(
						x.Body, func(n ast.Node) bool {
							switch call := n.(type) {
							case *ast.CallExpr:
								selExpr, ok := call.Fun.(*ast.SelectorExpr)
								if !ok {
									return true
								}

								ident, ok := selExpr.X.(*ast.Ident)
								if !ok {
									return true
								}

								if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
									pass.Reportf(
										n.Pos(), "prohibits the use of os.Exit in main",
									)
								}
							}
							return true
						},
					)
				}
				return true
			},
		)
	}

	return nil, nil
}
