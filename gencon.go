package gencon

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "gencon is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gencon",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	m := make(map[string]map[types.Type]bool)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CallExpr:
			var ident *ast.Ident
			switch fun := n.Fun.(type) {
			case *ast.IndexExpr:
				idt, ok := fun.X.(*ast.Ident)
				if !ok {
					// log.Fatalln("type assertion error: *ast.IndexExpr.fun.X to *ast.Ident")
					return
				}
				ident = idt
			case *ast.IndexListExpr:
				idt, ok := fun.X.(*ast.Ident)
				if !ok {
					// log.Fatalln("type assertion error: *ast.IndexListExpr.fun.X to *ast.Ident")
					return
				}
				ident = idt
			case *ast.Ident:
				ident = fun
			}
			instance, ok := pass.TypesInfo.Instances[ident]
			if !ok {
				return
			}
			id := pass.TypesInfo.ObjectOf(ident).Id()
			if m[id] == nil {
				m[id] = make(map[types.Type]bool)
			}

			typeArgs := instance.TypeArgs
			for i := 0; i < typeArgs.Len(); i++ {
				typ := typeArgs.At(i)
				if m[id][typ] {
					continue
				}
				ultyp := typ.Underlying()
				if !types.Identical(typ, ultyp) {
					m[id][ultyp] = true
					continue
				}
				m[id][typ] = false
			}
		}
	})

	return nil, nil
}
