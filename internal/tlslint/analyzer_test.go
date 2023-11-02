package tlslint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/go-tls-lint/internal/tlslint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}
	repoDir := filepath.Dir(filepath.Dir(wd))
	testdata := filepath.Join(repoDir, "testdata")

	analysistest.Run(
		t,
		testdata,
		tlslint.Analyzer,
		"basic",
	)
}
