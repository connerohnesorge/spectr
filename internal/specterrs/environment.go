package specterrs

// EditorNotSetError indicates the EDITOR environment variable is not set.
type EditorNotSetError struct {
	Operation string
}

func (*EditorNotSetError) Error() string {
	return "EDITOR environment variable not set"
}
