package commands

import (
	"fmt"
	"github.com/romnn/embedme/internal"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type EmbedFileCommand struct {
	Command
	Root      string
	Path      string
	StartLine int
	EndLine   int
	Base      string
}

func NewEmbedFileCommand(base string) *EmbedFileCommand {
	return &EmbedFileCommand{
		Path:      "",
		StartLine: 0,
		EndLine:   0,
		Base:      base,
	}
}

func (cmd *EmbedFileCommand) Lines() (int, int, bool) {
	valid := true
	if cmd.StartLine >= cmd.EndLine {
		valid = false
	}
	return cmd.StartLine, cmd.EndLine, valid
}

func (cmd *EmbedFileCommand) Output() ([]string, error) {
	path := filepath.Join(cmd.Base, cmd.Path)
	if err := internal.EnsureFile(path); err != nil {
		return nil, fmt.Errorf("failed to embed %s: %v", cmd.Path, err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// convert to lines
	newline := internal.DetectNewline(content)
	lines := internal.Lines(string(content), newline)

  // select lines
	if startLine, endLine, ok := cmd.Lines(); ok {
		lines = lines[startLine:endLine]
	}

  // properly indent
	minSpaces := internal.MinIndent(lines)
	if minSpaces > 0 {
		for i, line := range lines {
			lines[i] = line[minSpaces:]
		}
	}
	return lines, nil
}

var (
	embedPathRegex = regexp.MustCompile(`^\s*(?P<path>\S+?)(#L?(?P<start>\d+)-L?(?P<end>\d+))?\s*$`)
)

func (cmd *EmbedFileCommand) Parse(comment string) error {
	matches := internal.GetMatches(embedPathRegex, comment)
	if len(matches) < 1 {
		return fmt.Errorf("%s is not a valid file command", comment)
	}
	match := matches[0]
	cmd.Path = match["path"].Text
	cmd.StartLine, _ = strconv.Atoi(match["start"].Text)
	cmd.EndLine, _ = strconv.Atoi(match["end"].Text)
	return nil
}
