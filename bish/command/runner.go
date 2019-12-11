package command

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dalloriam/bish/bish/hooks"

	"github.com/dalloriam/bish/bish/state"
)

type CommandRequest struct {
	Context   *state.State
	UserInput string
	Hooks     []hooks.Hook

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *CommandRequest) Execute() error {
	// Remove the newline character.
	input := strings.TrimSuffix(c.UserInput, "\n")

	args, err := ParseArguments(input)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	planner := NewExecutionPlanner(c.Context, c.Hooks, args)
	cmd, err := planner.Command(true)
	if err != nil {
		return err
	}

	cmd.Bind(c.Stdin, c.Stdout, c.Stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		cmd.Signal(sig)
	}()

	_, err = cmd.Wait()
	return err
}
