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
	ctx := kong.Parse(
		cli,
		kong.Name("spectr"),
		kong.Description(
			"Validatable spec-driven development\n\n"+
				"Configuration: ~/.config/spectr/config.yaml "+
				"(run 'spectr config' to view)",
		),
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
