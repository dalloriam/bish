package bish

import (
	"fmt"
	"os"

	"github.com/dalloriam/bish/version"

	"github.com/dalloriam/bish/bish/command"
)

// Shell is the root bish shell struct.
type Shell struct {
	backend ShellBackend
}

// New initializes & returns a new shell instance with the provided IO writers.
func New(backend ShellBackend) *Shell {
	return &Shell{backend: backend}
}

func (s *Shell) prompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("wduss@phoenix [%s]\n> ", cwd)
}

func (s *Shell) err(err error) {
	s.backend.Stderr(fmt.Sprintf("%s\n", err.Error()))
}

// Start starts the main shell loop.
func (s *Shell) Start() {
	s.backend.Stdout(fmt.Sprintf("BiSH %s\n", version.VERSION))
	for {
		s.backend.Stdout(s.prompt())

		rawLine, err := s.backend.ReadLine()
		if err != nil {
			s.err(err)
			continue
		}

		if err := command.DoCommand(rawLine); err != nil {
			s.err(err)
			continue
		}
	}
}
