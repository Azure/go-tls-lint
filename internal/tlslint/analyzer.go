package tlslint

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const typeNameCryptoTLSConfig = "crypto/tls.Config"
const typeNameCryptoTLSConfigPointer = "*crypto/tls.Config"

var disableIgnoreComment bool

var Analyzer = &analysis.Analyzer{
	Name: "tls",
	Doc:  "Checks common TLS configuration mistakes.",
	URL:  "https://github.com/Azure/go-tls-lint",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	ResultType: reflect.TypeOf([]*Issue(nil)),
}

func init() {
	Analyzer.Flags.BoolVar(&disableIgnoreComment, "disable-ignore-comment", false, "disallow ignoring issues with a comment")
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

	isValidIssue := func(issue *Issue) bool {
		if issue.Severity == IssueSeverityError {
			return true
		}

		if disableIgnoreComment {
			return true
		}

		if hasIgnoreComment(commentMaps, issue.Node) {
			return false
		}

		return true
	}

	var issues []*Issue

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CompositeLit:
			issues = append(issues, filter(handleCompositeLit(pass, n), isValidIssue)...)
		case *ast.AssignStmt:
			issues = append(issues, filter(handleAssignStmt(pass, n), isValidIssue)...)
		}
	})

	for _, issue := range issues {
		pass.Reportf(issue.Node.Pos(), issue.String())
	}

	return issues, nil
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

func handleTLSConfigValue(
	n ast.Node,
	key ast.Expr,
	value ast.Expr,
) *Issue {
	if key == nil {
		return nil
	}

	keyIdent, ok := key.(*ast.Ident)
	if !ok {
		return nil
	}

	if keyIdent.Name == "" {
		return nil
	}

	issue := &Issue{
		Node: n,
	}

	if reason, exists := tlsConfigNamesBlockList[keyIdent.Name]; exists {
		issue.Severity = IssueSeverityError
		issue.Message = reason
	} else {
		issue.Severity = IssueSeverityWarning
		issue.Message = fmt.Sprintf("Unexpected TLS config settings %q", keyIdent.Name)
	}

	return issue
}

func handleCompositeLit(pass *analysis.Pass, n *ast.CompositeLit) []*Issue {
	if !isTLSConfigType(pass, n.Type) {
		return nil
	}

	var rv []*Issue
	for _, e := range n.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		issue := handleTLSConfigValue(e, kv.Key, kv.Value)
		if issue != nil {
			rv = append(rv, issue)
		}
	}

	return rv
}

func handleAssignStmt(pass *analysis.Pass, n *ast.AssignStmt) []*Issue {
	if len(n.Lhs) < 1 || len(n.Rhs) < 1 {
		return nil
	}
	lhs, ok := n.Lhs[0].(*ast.SelectorExpr) // TODO: handle multiple LHS
	if !ok {
		return nil
	}
	if !isTLSConfigType(pass, lhs.X) {
		return nil
	}

	issue := handleTLSConfigValue(n, lhs.Sel, n.Rhs[0])
	if issue != nil {
		return []*Issue{issue}
	}

	return nil
}
