package builtins

import (
	"os"

	"github.com/dalloriam/bish/bish/state"
)

// Name of the exit builtin.
const (
	exitName = "exit"
)

func init() {
	registry[exitName] = exit
}

func exit(*state.State, []string) (string, error) {
	os.Exit(0)
	return "", nil
}
