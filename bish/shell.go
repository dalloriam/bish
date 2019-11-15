package bish

import (
	"fmt"
	"os"

	"github.com/dalloriam/bish/bish/command"

	"github.com/dalloriam/bish/bish/config"
)

// Shell is the root bish shell struct.
type Shell struct {
	cfg config.IOConfig
}

// New initializes & returns a new shell instance with the provided IO writers.
func New(cfg config.IOConfig) *Shell {
	return &Shell{cfg: cfg}
}

func (s *Shell) prompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("wduss@phoenix [%s]\n> ", cwd)
}

func (s *Shell) stdout(a ...interface{}) {
	fmt.Fprint(s.cfg.Stdout, a...)
}

func (s *Shell) stderr(a ...interface{}) {
	fmt.Fprint(s.cfg.Stderr, a...)
}

func (s *Shell) err(err error) {
	s.stderr(fmt.Sprintf("error: %s", err.Error()))
}

// Start starts the main shell loop.
func (s *Shell) Start() {
	cmdReader := command.NewReader(s.cfg.Stdin)
	for {
		s.stdout(s.prompt())

		cmd, err := cmdReader.Read()
		if err != nil {
			s.err(err)
		}

		if err := cmd.Execute(s.cfg); err != nil {
			s.err(err)
		}
	}
}
