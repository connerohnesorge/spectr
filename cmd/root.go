package cmd

import (
	"github.com/connerohnesorge/spectr/internal/archive"
	kongcompletion "github.com/jotaen/kong-completion"
)

// CLI represents the root command structure for Kong
//
//nolint:revive // line-length-limit - struct tags cannot be split
type CLI struct {
	Init       InitCmd                   `cmd:"" help:"Initialize Spectr"`
	List       ListCmd                   `cmd:"" help:"List changes or specs"`
	Validate   ValidateCmd               `cmd:"" help:"Validate items"`
	Archive    archive.ArchiveCmd        `cmd:"" help:"Archive a change"`
	PR         PRCmd                     `cmd:"" help:"Create pull request from change"`
	View       ViewCmd                   `cmd:"" help:"Display dashboard"`
	Version    VersionCmd                `cmd:"" help:"Show version info"`
	Completion kongcompletion.Completion `cmd:"" help:"Generate completions"`
}
