package archive

import (
	"github.com/connerohnesorge/spectr/internal/list"
)

// newListerForArchive creates a lister for the archive package
func newListerForArchive(projectPath string) *list.Lister {
	return list.NewLister(projectPath)
}

// runInteractiveArchiveForArchiver wraps the list package's
// interactive archive function
func runInteractiveArchiveForArchiver(
	changes []list.ChangeInfo,
	projectPath string,
) (string, error) {
	return list.RunInteractiveArchive(changes, projectPath)
}
