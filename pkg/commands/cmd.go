package commands

import (
	"fmt"
	"github.com/romnn/embedme/internal"
	"regexp"
  "os/exec"
)

type EmbedCommandOutputCommand struct {
	Command
	Cmd string
	Cwd string
}

func NewEmbedCommandOutputCommand(cwd string) *EmbedCommandOutputCommand {
	return &EmbedCommandOutputCommand{
		Cmd: "",
		Cwd: cwd,
	}
}

var (
	embedCommandOutputRegex = regexp.MustCompile(`^\s*\$\s*(?P<command>[\s\S]+?)\s*$`)
)

func (cmd *EmbedCommandOutputCommand) Parse(comment string) error {
	matches := internal.GetMatches(embedCommandOutputRegex, comment)
	if len(matches) < 1 {
		return fmt.Errorf("%s is not a valid command", comment)
	}
	match := matches[0]
	cmd.Cmd = match["command"].Text
	return nil
}

func (cmd *EmbedCommandOutputCommand) Output() ([]string, error) {
  execCmd := exec.Command("sh", "-c", cmd.Cmd)
  // set working directory of the command
  execCmd.Dir = cmd.Cwd
  output, err := execCmd.CombinedOutput()
  if err != nil {
    return nil, err
  }
	newline := internal.DetectNewline(output)
	lines := internal.Lines(string(output), newline)
  return lines, nil
}
