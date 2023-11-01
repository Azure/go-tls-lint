package main

import "github.com/alecthomas/kong"

func main() {
	ctx := kong.Parse(&Probe{})

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
