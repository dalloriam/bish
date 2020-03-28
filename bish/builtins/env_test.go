package builtins_test

import (
	"os"

	"github.com/dalloriam/bish/bish/state"
)

func validateEnv(ctx *state.State, args []string) bool {
	if len(args) != 3 {
		return false
	}

	key := args[1]
	val := args[2]

	v := os.Getenv(key)

	return val == v
}
