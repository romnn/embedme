package commands

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/romnn/embedme/internal"
	"github.com/spf13/afero"
)

// EmbedCommandOutputCommand ...
type EmbedCommandOutputCommand struct {
	Command
	Cmd        string
	WorkingDir string
	FS         afero.Fs
}

// NewEmbedCommandOutputCommand ...
func NewEmbedCommandOutputCommand(fs afero.Fs, cwd string) *EmbedCommandOutputCommand {
	return &EmbedCommandOutputCommand{
		Cmd:        "",
		WorkingDir: cwd,
		FS:         fs,
	}
}

var (
	embedCommandOutputRegex = regexp.MustCompile(`^\s*\$\s*(?P<command>[\s\S]+?)\s*$`)
)

// Parse ...
func (cmd *EmbedCommandOutputCommand) Parse(comment string) error {
	matches := internal.GetMatches(embedCommandOutputRegex, comment)
	if len(matches) < 1 {
		return fmt.Errorf("%s is not a valid command", comment)
	}
	match := matches[0]
	cmd.Cmd = match["command"].Text
	return nil
}

// Output ...
func (cmd *EmbedCommandOutputCommand) Output() ([]string, error) {
	execCmd := exec.Command("sh", "-c", cmd.Cmd)
	// set working directory of the command
	execCmd.Dir = cmd.WorkingDir
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := internal.Lines(string(output))
	return lines, nil
}
