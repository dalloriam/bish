package command

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type CommandRequest struct {
	Context   ShellContext
	UserInput string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func DoCommand(req CommandRequest) error {
	// Remove the newline character.
	input := strings.TrimSuffix(req.UserInput, "\n")

	args, err := ParseArguments(input)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	planner := NewExecutionPlanner(req.Context, args)
	cmd, err := planner.Command(true)
	if err != nil {
		return err
	}

	cmd.Bind(req.Stdin, req.Stdout, req.Stderr)

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
