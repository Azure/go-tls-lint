package main

import (
	"flag"

	"github.com/Azure/go-tls-lint/internal/version"
)

func init() {
	versionFlag := version.Flag{Prog: "go-tls-lint"}

	// NOTE: the V usage is to align with the go tool.
	flag.Var(versionFlag, "V", "print version and exit")
	flag.Var(versionFlag, "version", "print version and exit")
}
