package command

import (
	"io"
	"os"
)

// The Command interface represents a comannd/subcommand that can be executed by
// the shell.
type Command interface {
	Argument

	Bind(stdin io.Reader, stdout, stderr io.Writer)

	Start() error
	Wait() (string, error)
	Signal(os.Signal)
}
