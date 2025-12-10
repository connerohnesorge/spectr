package specterrs

// EmptyRemoteURLError indicates an empty remote URL was encountered.
type EmptyRemoteURLError struct{}

func (*EmptyRemoteURLError) Error() string {
	return "empty remote URL"
}

// BranchNameRequiredError indicates a branch name is required but
// was not provided.
type BranchNameRequiredError struct{}

func (*BranchNameRequiredError) Error() string {
	return "branch name is required"
}

// BaseBranchRequiredError indicates a base branch is required but
// was not provided.
type BaseBranchRequiredError struct{}

func (*BaseBranchRequiredError) Error() string {
	return "base branch is required"
}

// NotInGitRepositoryError indicates the path is not within a git repository.
type NotInGitRepositoryError struct {
	Path string
}

func (*NotInGitRepositoryError) Error() string {
	return "not a git repository"
}

// BaseBranchNotFoundError indicates the base branch could not be determined.
type BaseBranchNotFoundError struct {
	BranchName string
}

func (*BaseBranchNotFoundError) Error() string {
	return "could not determine base branch, please specify with --base"
}
