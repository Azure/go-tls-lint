package tlslint

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// lazyCommentMaps provides lazy initialization of comment maps from a set of files.
type lazyCommentMaps struct {
	materialized map[*ast.File]ast.CommentMap
	fset         *token.FileSet
	files        []*ast.File
}

func newLazyCommentMaps(pass *analysis.Pass) *lazyCommentMaps {
	return &lazyCommentMaps{
		materialized: make(map[*ast.File]ast.CommentMap),
		fset:         pass.Fset,
		files:        pass.Files,
	}
}

// GetCommentForNode returns the comment groups associated with the given node.
func (m *lazyCommentMaps) GetCommentForNode(n ast.Node) ([]*ast.CommentGroup, bool) {
	for _, materialized := range m.materialized {
		if c, ok := materialized[n]; ok {
			return c, true
		}
	}

	for _, f := range m.files {
		if _, ok := m.materialized[f]; ok {
			continue
		}

		m.materialized[f] = ast.NewCommentMap(m.fset, f, f.Comments)
		if c, ok := m.materialized[f][n]; ok {
			return c, true
		}
	}

	return nil, false
}

var ignoreComment = regexp.MustCompile("^go-tls-lint:ignore (.*)$")

func hasIgnoreComment(commentMaps *lazyCommentMaps, node ast.Node) bool {
	comments, hasComments := commentMaps.GetCommentForNode(node)
	if !hasComments {
		return false
	}

	for _, comment := range comments {
		commentText := strings.TrimSpace(comment.Text()) // NOTE: no // prefix
		if ignoreComment.MatchString(commentText) {
			return true
		}
	}

	return false
}
