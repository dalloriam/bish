package builtins

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/dalloriam/bish/bish/state"
	"github.com/dalloriam/bish/version"
)

// Name of the version builtin.
const (
	versionName = "version"
)

func init() {
	registry[versionName] = getVersion
}

func getVersion(*state.State, []string) (string, error) {
	bish := aurora.Green("BiSH")
	vPrompt := aurora.Blue("Version: ")
	gPrompt := aurora.Blue("Git Commit: ")
	return fmt.Sprintf("%s\n%s%s\n%s%s\n", bish, vPrompt, version.VERSION, gPrompt, version.GITCOMMIT), nil
}
