package bish

import (
	"fmt"
	"os"

	"github.com/dalloriam/bish/bish/command"
	. "github.com/logrusorgru/aurora"
)

// Shell is the root bish shell struct.
type Shell struct {
	backend ShellBackend
	ctx     *ContextStore
}

// New initializes & returns a new shell instance with the provided IO writers.
func New(backend ShellBackend) *Shell {
	return &Shell{backend: backend, ctx: NewContext()}
}

func (s *Shell) prompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("wduss@phoenix [%s]\n> ", cwd)
}

func (s *Shell) err(err error) {
	_, err = s.backend.Stderr().Write([]byte(fmt.Sprintf("! %s\n", Red(err.Error()))))
	return
}

func (s *Shell) stdout(d string) (err error) {
	_, err = s.backend.Stdout().Write([]byte(d))
	return
}

// Start starts the main shell loop.
func (s *Shell) Start() {
	for {
		if err := s.stdout(s.prompt()); err != nil {
			s.err(err)
			continue
		}

		rawLine, err := s.backend.ReadLine()
		if err != nil {
			s.err(err)
			continue
		}

		if err := command.DoCommand(command.CommandRequest{
			Context:   s.ctx,
			UserInput: rawLine,
			Stdin:     s.backend.Stdin(),
			Stdout:    s.backend.Stdout(),
			Stderr:    s.backend.Stderr(),
		}); err != nil {
			s.err(err)
			continue
		}
	}
}
