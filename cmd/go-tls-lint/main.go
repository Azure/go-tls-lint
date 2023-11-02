package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Azure/go-tls-lint/internal/tlslint"
)

func main() {
	singlechecker.Main(tlslint.Analyzer)
}
