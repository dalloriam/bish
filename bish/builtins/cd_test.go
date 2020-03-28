package builtins_test

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/dalloriam/bish/bish/state"
)

func validateCd(ctx *state.State, args []string) bool {
	if len(args) > 2 {
		return false
	}

	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	if len(args) == 1 {
		// Cd to home dir.

		user, err := user.Current()
		if err != nil {
			return false
		}

		return wd == user.HomeDir
	} else {
		wd, err = filepath.Abs(wd)
		if err != nil {
			return false
		}
		return wd == args[1]
	}
}
