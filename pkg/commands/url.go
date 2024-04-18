package commands

import "net/url"

type EmbedURLCommand struct {
	Command
	// TODO: add lines here? for text based files
	URL url.URL
}

func NewEmbedURLCommand(url url.URL) *EmbedURLCommand {
	return &EmbedURLCommand{
		URL: url,
	}
}

func (cmd *EmbedURLCommand) Output() ([]string, error) {
	return []string{}, nil
}

func (cmd *EmbedURLCommand) Parse(comment string) error {
	return nil
}
