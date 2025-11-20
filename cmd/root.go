package cmd

import (
	"github.com/connerohnesorge/spectr/internal/archive"
)

// CLI represents the root command structure for Kong
type CLI struct {
	Init     InitCmd            `cmd:"" help:"Initialize Spectr in a project"`
	List     ListCmd            `cmd:"" help:"List changes or specifications"`
	Validate ValidateCmd        `cmd:"" help:"Validate changes or specs"`
	Archive  archive.ArchiveCmd `cmd:"" help:"Archive a completed change"`
	View     ViewCmd            `cmd:"" help:"Display project dashboard"`
}
