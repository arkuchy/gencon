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
	// {objectId: {type: bool}}
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

	nodeFilter = []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			id := pass.TypesInfo.ObjectOf(n.Name).Id()
			tp := n.Type.TypeParams
			if tp == nil {
				return
			}
			tps := tp.List

			// TODO
			if len(tps) > 1 {
				return
			}
			for _, f := range tps {
				// TODO
				if len(f.Names) > 1 {
					return
				}
				tv := pass.TypesInfo.Types[f.Type]
				// FIX ME: do not compare with Type.String()
				if tv.Type.String() == "any" {
					pass.Reportf(f.Pos(), "change any to %v", m[id])
				}
				// for _, ident := range f.Names {
				// 	object := pass.TypesInfo.ObjectOf(ident)
				// 	fmt.Println(object.Name())
				// 	fmt.Println(object.Type())
				// 	fmt.Println(m[id])
				// }

			}
		}
	})

	return nil, nil
}
