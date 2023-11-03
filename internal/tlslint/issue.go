package tlslint

import (
	"fmt"
	"go/ast"
)

// IssueSeverity is the severity of an issue.
type IssueSeverity string

const (
	// IssueSeverityError is an error issue, which cannot be ignored by code comment.
	IssueSeverityError IssueSeverity = "ERROR"
	// IssueSeverityWarning is a warning issue, which can be ignored by code comment.
	IssueSeverityWarning IssueSeverity = "WARNING"
)

// Issue is a linting issue.
type Issue struct {
	// Severity - the severity of the issue.
	Severity IssueSeverity
	// Message - the message to display to the user.
	Message string
	// Node - the node that triggered the issue.
	Node ast.Node
}

func (i *Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Severity, i.Message)
}

// tlsConfigNamesBlockList contains a list of TLS config names that should not be used.
// Settings any of these will result in an error issue.
var tlsConfigNamesBlockList = map[string]string{
	"MinVersion":   "Go 1.18 and onward has good defaults for MinVersion, no need to set it",
	"MaxVersion":   "Don't pin MaxVersion which disables future TLS versions",
	"CipherSuites": "Overriding CipherSuites is dangerous, don't do it unless you know what you are doing",
}

// tlsConfigNamesWarnList contains a list of TLS config names that should avoid being used.
// Settings any of these will result in a warning issue.
var tlsConfigNamesWarnList = map[string]string{
	"Renegotiation": "Renegotiation is not supported in TLS 1.3",
}

func filter[T any](xs []T, f func(T) bool) []T {
	var ys []T
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}
