package builtins

import "os"

// Name of the exit builtin.
const (
	ExitName = "exit"
)

// Exit exits the shell.
func Exit() error {
	os.Exit(0)
	return nil
}
