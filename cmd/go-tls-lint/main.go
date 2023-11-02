package main

import (
	"github.com/Azure/go-tls-lint/internal/tlslint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(tlslint.Analyzer)
}
