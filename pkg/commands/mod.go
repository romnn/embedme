package commands

// Command is a generic embedme command
type Command interface {
	// Output gets the output lines that should be embedded
	Output() ([]string, error)
	// Parse parses a comment as this command
	//
	// If the comment cannot be parsed, an error is returned.
	Parse(comment string) error
}
