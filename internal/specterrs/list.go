package specterrs

import "fmt"

// IncompatibleFlagsError indicates two flags cannot be used together.
type IncompatibleFlagsError struct {
	Flag1 string
	Flag2 string
}

func (e *IncompatibleFlagsError) Error() string {
	return fmt.Sprintf(
		"cannot use %s with %s",
		e.Flag1,
		e.Flag2,
	)
}

// RequiresFlagError indicates a flag requires another flag to be set.
type RequiresFlagError struct {
	Flag         string
	RequiredFlag string
}

func (e *RequiresFlagError) Error() string {
	return fmt.Sprintf(
		"%s requires %s",
		e.Flag,
		e.RequiredFlag,
	)
}
