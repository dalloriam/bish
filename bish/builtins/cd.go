package builtins

import (
	"errors"
	"os"
	"os/user"

	"github.com/dalloriam/bish/bish/state"
)

const (
	cdName = "cd"
)

func init() {
	registry[cdName] = changeDirectory
}

func changeDirectory(ctx *state.State, args []string) (string, error) {
	if len(args) == 0 {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return "", os.Chdir(usr.HomeDir)
	} else if len(args) == 1 {
		return "", os.Chdir(args[0])
	}
	return "", errors.New("Too many arguments for cd command")
}
