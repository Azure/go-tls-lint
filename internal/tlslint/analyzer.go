package tlslint

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const typeNameCryptoTLSConfig = "crypto/tls.Config"
const typeNameCryptoTLSConfigPointer = "*crypto/tls.Config"

var tlsConfigNamesBlockList = map[string]string{
	"MinVersion": "ERROR: go 1.18 and onward has good defaults for MinVersion, no need to set it",
	"MaxVersion": "ERROR: don't pin MaxVersion which disables future TLS versions",
}

var Analyzer = &analysis.Analyzer{
	Name: "tls",
	Doc:  "Checks TLS configurations",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}


func run(pass *analysis.Pass) (any, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		// a := &tls.Config{}
		(*ast.CompositeLit)(nil),

		// a.MinVersion = tls.VersionTLS12
		(*ast.AssignStmt)(nil),
	}

	commentMaps := newLazyCommentMaps(pass)

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CompositeLit:
			handleCompositeLit(pass, commentMaps, n)
		case *ast.AssignStmt:
			handleAssignStmt(pass, commentMaps, n)
		}
	})

	return nil, nil
}

func isTLSConfigType(pass *analysis.Pass, n ast.Expr) bool {
	if n == nil {
		return false
	}

	actualType := pass.TypesInfo.TypeOf(n)
	if actualType == nil {
		return false
	}

	switch actualType := actualType.(type) {
	case *types.Named:
		return actualType.String() == typeNameCryptoTLSConfig
	case *types.Pointer:
		return actualType.String() == typeNameCryptoTLSConfigPointer
	default:
		return false
	}
}

func processTLSConfigValue(pass *analysis.Pass, key ast.Expr, value ast.Expr) string {
	if key == nil {
		return ""
	}

	keyIdent, ok := key.(*ast.Ident)
	if !ok {
		return ""
	}

	if keyIdent.Name == "" {
		return ""
	}

	if reason, exists := tlsConfigNamesBlockList[keyIdent.Name]; exists {
		// something
		return reason
	}

	return fmt.Sprintf("WARN: unexpected TLS config settings %q", keyIdent.Name)
}

func handleCompositeLit(pass *analysis.Pass, commentMaps *lazyCommentMaps, n *ast.CompositeLit) {
	if !isTLSConfigType(pass, n.Type) {
		return
	}

	for _, e := range n.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		report := processTLSConfigValue(pass, kv.Key, kv.Value)
		if report == "" {
			continue
		}

		if hasIgnoreComment(commentMaps, kv) {
			continue
		}

		pass.Reportf(kv.Pos(), report)
	}
}

func handleAssignStmt(pass *analysis.Pass, commentMaps *lazyCommentMaps, n *ast.AssignStmt) {
	if len(n.Lhs) < 1 || len(n.Rhs) < 1 {
		return
	}
	lhs, ok := n.Lhs[0].(*ast.SelectorExpr) // TODO: handle multiple LHS
	if !ok {
		return
	}
	if !isTLSConfigType(pass, lhs.X) {
		return
	}

	report := processTLSConfigValue(pass, lhs.Sel, n.Rhs[0])
	if report == "" {
		return
	}

	if hasIgnoreComment(commentMaps, n) {
		return
	}

	pass.Reportf(n.Pos(), report)
}
