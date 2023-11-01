package main

import "github.com/alecthomas/kong"

func main() {
	cli := &Probe{}
	ctx := kong.Parse(cli)

	ctx.BindTo(&cliProvider{Verbose: cli.Verbose}, (*Provider)(nil))

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
