package tlslint

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const typeNameCryptoTLSConfig = "crypto/tls.Config" 

var tlsConfigNamesBlockList = map[string]string{
	"MinVersion": "Error: go 1.18 and onward has good defaults for MinVersion, no need to set it",
	"MaxVersion": "Error: don't pin MaxVersion which disables future TLS versions",
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

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CompositeLit:
			handleCompositeLit(pass, n)
		case *ast.AssignStmt:
			handleAssignStmt(pass, n)
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
	return  actualType.String() == typeNameCryptoTLSConfig
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

func handleCompositeLit(pass *analysis.Pass, n *ast.CompositeLit) {
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

		pass.Reportf(kv.Pos(), report)
	}
}

func handleAssignStmt(pass *analysis.Pass, n *ast.AssignStmt) {
}