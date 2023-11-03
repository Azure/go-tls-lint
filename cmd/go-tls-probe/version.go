package main

import (
	"fmt"

	"github.com/Azure/go-tls-lint/internal/version"
	"github.com/alecthomas/kong"
)

type VersionFlag bool

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (VersionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Println(version.GetFull("go-tls-probe"))
	app.Exit(0)
	return nil
}
