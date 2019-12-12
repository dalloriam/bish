package builtins

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/dalloriam/bish/version"
)

// Name of the version builtin.
const (
	VersionName = "version"
)

func Version() (string, error) {
	bish := aurora.Green("BiSH")
	vPrompt := aurora.Blue("Version: ")
	gPrompt := aurora.Blue("Git Commit: ")
	return fmt.Sprintf("%s\n%s%s\n%s%s\n", bish, vPrompt, version.VERSION, gPrompt, version.GITCOMMIT), nil
}
