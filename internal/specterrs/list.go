package specterrs

import "fmt"

// IncompatibleFlagsError indicates two flags cannot be used together.
type IncompatibleFlagsError struct {
	Flag1 string
	Flag2 string
}

func (e *IncompatibleFlagsError) Error() string {
	return fmt.Sprintf("cannot use %s with %s", e.Flag1, e.Flag2)
}
