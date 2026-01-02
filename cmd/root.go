package cmd

import (
	"github.com/connerohnesorge/spectr/internal/archive"
	kongcompletion "github.com/jotaen/kong-completion"
)

// CLI represents the root command structure for Kong
type CLI struct {
	Init       InitCmd                   `cmd:"" help:"Initialize Spectr"`                        //nolint:lll,revive // Kong struct tag with alignment
	List       ListCmd                   `cmd:"" help:"List items"                  aliases:"ls"` //nolint:lll,revive // Kong struct tag with alignment
	Validate   ValidateCmd               `cmd:"" help:"Validate items"`                           //nolint:lll,revive // Kong struct tag with alignment
	Accept     AcceptCmd                 `cmd:"" help:"Accept tasks.md"`                          //nolint:lll,revive // Kong struct tag with alignment
	Archive    archive.ArchiveCmd        `cmd:"" help:"Archive a change"`                         //nolint:lll,revive // Kong struct tag with alignment
	PR         PRCmd                     `cmd:"" help:"Create pull requests"`                     //nolint:lll,revive // Kong struct tag with alignment
	Tasks      TasksCmd                  `cmd:"" help:"Show task status"`                         //nolint:lll,revive // Kong struct tag with alignment
	View       ViewCmd                   `cmd:"" help:"Display dashboard"`                        //nolint:lll,revive // Kong struct tag with alignment
	Version    VersionCmd                `cmd:"" help:"Show version info"`                        //nolint:lll,revive // Kong struct tag with alignment
	Completion kongcompletion.Completion `cmd:"" help:"Generate completions"`                     //nolint:lll,revive // Kong struct tag with alignment
}
