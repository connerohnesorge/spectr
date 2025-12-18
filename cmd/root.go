package cmd

import (
	"github.com/connerohnesorge/spectr/internal/archive"
	kongcompletion "github.com/jotaen/kong-completion"
)

// CLI represents the root command structure for Kong
type CLI struct {
	Init       InitCmd                   `cmd:"" help:"Initialize Spectr"`                        //nolint:lll,revive
	List       ListCmd                   `cmd:"" help:"List items"                  aliases:"ls"` //nolint:lll,revive
	Validate   ValidateCmd               `cmd:"" help:"Validate items"`                           //nolint:lll,revive
	Accept     AcceptCmd                 `cmd:"" help:"Accept tasks.md"`                          //nolint:lll,revive
	Archive    archive.ArchiveCmd        `cmd:"" help:"Archive a change"`                         //nolint:lll,revive
	Track      TrackCmd                  `cmd:"" help:"Auto-commit on task changes"`              //nolint:lll,revive
	PR         PRCmd                     `cmd:"" help:"Create pull requests"`                     //nolint:lll,revive
	View       ViewCmd                   `cmd:"" help:"Display dashboard"`                        //nolint:lll,revive
	Version    VersionCmd                `cmd:"" help:"Show version info"`                        //nolint:lll,revive
	Completion kongcompletion.Completion `cmd:"" help:"Generate completions"`                     //nolint:lll,revive
}
