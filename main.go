/*
Copyright © 2025 Conner Ohnesorge
*/
package main

import (
	"github.com/alecthomas/kong"
	"github.com/connerohnesorge/spectr/cmd"
	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/theme"
)

func main() {
	cli := &cmd.CLI{}
	ctx := kong.Parse(cli,
		kong.Name("spectr"),
		kong.Description("Validatable spec-driven development"),
		kong.UsageOnError(),
	)

	// Load config and apply theme
	cfg, err := config.Load()
	if err == nil {
		_ = theme.Load(cfg.Theme)
	}
	// Ignore errors - theme will default to "default" if config not found

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
