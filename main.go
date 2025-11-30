/*
Copyright © 2025 Conner Ohnesorge
*/
package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/connerohnesorge/spectr/cmd"
	kongcompletion "github.com/jotaen/kong-completion"
)

func main() {
	cli := &cmd.CLI{}
	app := kong.Must(cli,
		kong.Name("spectr"),
		kong.Description("Validatable spec-driven development"),
		kong.UsageOnError(),
	)

	// Register shell completion with custom predictors
	kongcompletion.Register(app,
		kongcompletion.WithPredictor("changeID", cmd.PredictChangeIDs()),
		kongcompletion.WithPredictor("specID", cmd.PredictSpecIDs()),
		kongcompletion.WithPredictor("itemType", cmd.PredictItemTypes()),
		kongcompletion.WithPredictor("item", cmd.PredictItems()),
	)

	ctx, err := app.Parse(os.Args[1:])
	app.FatalIfErrorf(err)
	err = ctx.Run()
	app.FatalIfErrorf(err)
}
