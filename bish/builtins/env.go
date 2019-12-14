package builtins

import (
	"errors"
	"os"

	"github.com/dalloriam/bish/bish/state"
)

// Name of the env command.
const (
	setEnvName = "set"
)

func init() {
	registry[setEnvName] = setEnv
}

func setEnv(_ *state.State, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("invalid setenv syntax")
	}

	key := args[0]
	val := args[1]
	return "", os.Setenv(key, val)
}
