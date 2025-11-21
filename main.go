/*
Copyright © 2025 Conner Ohnesorge
*/
package main

import (
	"github.com/alecthomas/kong"
	"github.com/connerohnesorge/spectr/cmd"
)

func main() {
	cli := &cmd.CLI{}
	ctx := kong.Parse(cli,
		kong.Name("spectr"),
		kong.Description("Validatable spec-driven development"),
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
