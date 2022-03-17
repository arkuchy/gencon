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

	nodeCallExprFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	m := make(map[*types.TypeParam]map[types.Type]bool)
	inspect.Preorder(nodeCallExprFilter, func(n ast.Node) {
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
			obj := pass.TypesInfo.ObjectOf(ident)
			sig, ok := obj.Type().(*types.Signature)
			if !ok {
				return
			}
			typeParams := sig.TypeParams()
			typeArgs := instance.TypeArgs

			if typeParams.Len() != typeArgs.Len() {
				return
			}
			for i := 0; i < typeArgs.Len(); i++ {
				typp := typeParams.At(i)
				typa := typeArgs.At(i)
				if m[typp] == nil {
					m[typp] = make(map[types.Type]bool)
				}
				if m[typp][typa] {
					continue
				}
				ultypa := typa.Underlying()
				if !types.Identical(typa, ultypa) {
					m[typp][ultypa] = true
					continue
				}
				m[typp][typa] = false
			}
		}
	})

	nodeFuncDeclFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	anyobj := types.Universe.Lookup("any")
	inspect.Preorder(nodeFuncDeclFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			sig, ok := pass.TypesInfo.TypeOf(n.Name).(*types.Signature)
			if !ok {
				return
			}
			typeParams := sig.TypeParams()
			tp := n.Type.TypeParams
			if tp == nil {
				return
			}
			fieldList := tp.List

			// we can check whether constraint is "any" or not like the following:
			//
			// for i := 0; i < typeParams.Len(); i++ {
			// 	typeParam := typeParams.At(i)
			// 	constr := typeParam.Constraint()
			// 	typInterface, _ := constr.(*types.Interface)
			// 	if typInterface.Empty() {
			// 		pass.Reportf()
			// 	}
			// }
			//
			// but it can detect empty interface(interface{}) etc.
			// so we use idx
			idx := 0
			for _, field := range fieldList {
				for _, name := range field.Names {
					typp := typeParams.At(idx)
					idx += 1
					idt, ok := field.Type.(*ast.Ident)
					if !ok {
						continue
					}
					obj := pass.TypesInfo.ObjectOf(idt)
					if obj == anyobj {
						pass.Reportf(name.Pos(), "change constraint of %s from any to %s", name, CreateUnion(m[typp]))
					}
				}
			}
		}
	})

	return nil, nil
}

// CreateUnion creates *types.Union from map.
func CreateUnion(m map[types.Type]bool) *types.Union {
	var terms []*types.Term
	for t, b := range m {
		terms = append(terms, types.NewTerm(b, t))
	}
	return types.NewUnion(terms)
}
