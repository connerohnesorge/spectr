package archive

import (
	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/list"
)

// newListerForArchiveWithConfig creates a lister for the archive package
// using config.
// Note: list package now supports config via NewListerWithConfig.
func newListerForArchiveWithConfig(cfg *config.Config) *list.Lister {
	return list.NewListerWithConfig(cfg)
}

// runInteractiveArchiveForArchiverWithConfig wraps the list package's
// interactive archive function with config.
// Note: list package uses cfg.ProjectRoot for interactive operations.
func runInteractiveArchiveForArchiverWithConfig(
	changes []list.ChangeInfo,
	cfg *config.Config,
) (string, error) {
	return list.RunInteractiveArchive(changes, cfg.ProjectRoot)
}
