package builtins

import (
	"errors"
	"os"
)

const (
	SetEnvName = "set"
)

func SetEnv(args []string) error {
	if len(args) != 2 {
		return errors.New("invalid setenv syntax")
	}

	key := args[0]
	val := args[1]
	return os.Setenv(key, val)
}
