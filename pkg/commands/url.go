package commands

import "net/url"

// EmbedURLCommand ...
type EmbedURLCommand struct {
	Command
	URL url.URL
}

// NewEmbedURLCommand ...
func NewEmbedURLCommand(url url.URL) *EmbedURLCommand {
	return &EmbedURLCommand{
		URL: url,
	}
}

// Output ...
func (cmd *EmbedURLCommand) Output() ([]string, error) {
	return []string{}, nil
}

// Parse ...
func (cmd *EmbedURLCommand) Parse(comment string) error {
	return nil
}
