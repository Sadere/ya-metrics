package linters

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Возвращает анализатор, который проверяет вызов стандартной функции os.Exit внутри пакета main в функции main
func NewExitAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "checkosexit",
		Doc:  "check if os.Exit is called inside main function of main package",
		Run:  NewExitRun(),
	}
}

func NewExitRun() func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		var osExitType types.Type

		for _, pkg := range pass.Pkg.Imports() {
			if pkg.Name() == "os" {
				exitObj := pkg.Scope().Lookup("Exit")

				if exitObj != nil {
					osExitType = exitObj.Type()
				}
			}
		}

		// os не импортирован
		if osExitType == nil {
			return nil, nil
		}

		// обходим все файлы
		for _, file := range pass.Files {
			// только main пакет
			if file.Name.Name != "main" {
				continue
			}

			// функцией ast.Inspect проходим по всем узлам AST
			ast.Inspect(file, func(node ast.Node) bool {
				// только функция main
				if f, ok := node.(*ast.FuncDecl); ok {
					if f.Name.Name != "main" {
						return false
					}
				}

				// проверяем вызов функции
				if callExpr, ok := node.(*ast.CallExpr); ok {
					callExprFunType := pass.TypesInfo.TypeOf(callExpr.Fun)

					if callExprFunType == osExitType {
						pass.Report(analysis.Diagnostic{
							Pos:     node.Pos(),
							End:     node.End(),
							Message: "usage of os.Exit inside main function of main package is prohibited",
						})
					}
				}

				return true
			})
		}

		return nil, nil
	}
}
