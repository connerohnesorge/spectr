package archive

import (
	"github.com/connerohnesorge/spectr/internal/list"
)

// newListerForArchive creates a lister for the archive package
func newListerForArchive(projectPath string) *list.Lister {
	return list.NewLister(projectPath)
}

// runInteractiveSelectChangeForArchiver wraps the list package's
// interactive change selection function.
func runInteractiveSelectChangeForArchiver(
	changes []list.ChangeInfo,
	projectPath string,
) (string, error) {
	return list.RunInteractiveSelectChange(changes, projectPath)
}
