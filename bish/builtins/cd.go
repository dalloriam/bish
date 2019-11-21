package builtins

import (
	"errors"
	"os"
	"os/user"
)

const CdName = "cd"

// ChangeDirectory changes the directory.
func ChangeDirectory(args []string) error {
	if len(args) == 0 {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		return os.Chdir(usr.HomeDir)
	} else if len(args) == 1 {
		return os.Chdir(args[0])
	}
	return errors.New("Too many arguments for cd command")
}
