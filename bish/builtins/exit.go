package builtins

import "os"

const ExitName = "exit"

// Exit exits the shell.
func Exit() error {
	os.Exit(0)
	return nil
}
