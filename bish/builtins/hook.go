package builtins

import (
	"errors"
	"io/ioutil"

	"github.com/dalloriam/bish/bish/script"
	"github.com/dalloriam/bish/bish/state"
)

// Name of the hook builtin.
const (
	hookName = "hook"
)

func init() {
	registry[hookName] = hook
}

func hook(ctx *state.State, args []string) (string, error) {
	if len(args) == 2 {
		// Setting a hook.
		hookSrc, err := ioutil.ReadFile(args[1])
		if err != nil {
			return "", err
		}
		hookObject, err := script.NewHook(args[0], string(hookSrc), ctx.ScriptRuntime())
		if err != nil {
			return "", err
		}
		ctx.AddHook(args[0], hookObject)
		return "", nil
	}
	return "", errors.New("invalid syntax")
}
