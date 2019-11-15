package command

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dalloriam/bish/bish/builtins"
)

// Command represents a command to execute.
type Command struct {
	Cmd       string
	Arguments []string
}

func ParseCommand(input string) (*Command, error) {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input to separate the command and the arguments.
	// TODO: Fancier argument parsing.
	args := strings.Split(input, " ")

	return &Command{Cmd: args[0], Arguments: args[1:]}, nil
}

func (c *Command) nativeExec() error {
	// Pass the program and the arguments separately.
	cmd := exec.Command(c.Cmd, c.Arguments...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}

// Execute executes the command.
func (c *Command) Execute() error {
	switch c.Cmd {
	case builtins.CdName:
		return builtins.ChangeDirectory(c.Arguments[0])
	case builtins.ExitName:
		return builtins.Exit()
	default:
		return c.nativeExec()
	}
}
