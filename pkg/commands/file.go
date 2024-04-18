package commands

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/romnn/embedme/internal"
	"github.com/romnn/embedme/pkg/fs"
	"github.com/spf13/afero"
)

type EmbedFileCommand struct {
	Command
	Root      string
	Path      string
	StartLine int
	EndLine   int
	BaseDirs  []string
	FS        afero.Fs
}

func NewEmbedFileCommand(fs afero.Fs, baseDirs ...string) *EmbedFileCommand {
	return &EmbedFileCommand{
		Path:      "",
		StartLine: 0,
		EndLine:   0,
		BaseDirs:  baseDirs,
		FS:        fs,
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
	candidatePaths := []string{}
	for _, base := range cmd.BaseDirs {
		candidatePaths = append(candidatePaths, filepath.Join(base, cmd.Path))
	}
	var existingPath string
	for _, path := range candidatePaths {
		if err := fs.EnsureFile(cmd.FS, path); err == nil {
			existingPath = path
		}
	}

	if existingPath == "" {
		return nil, fmt.Errorf(
			"failed to embed: neither of %v exists",
			candidatePaths,
		)
	}

	content, err := afero.ReadFile(cmd.FS, existingPath)
	if err != nil {
		return nil, err
	}

	// convert to lines
	// newline := internal.DetectNewline(content)
	// lines := internal.Lines(string(content), newline)
	lines := internal.Lines(string(content))

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
