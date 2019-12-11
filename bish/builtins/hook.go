package builtins

import (
	"errors"
	"io/ioutil"

	"github.com/dalloriam/bish/bish/hooks"

	"github.com/dalloriam/bish/bish/state"
)

const (
	HookName = "hook"
)

func Hook(ctx *state.State, args []string) error {
	if len(args) == 2 {
		// Setting a hook.
		hookSrc, err := ioutil.ReadFile(args[1])
		if err != nil {
			return err
		}
		hookObject, err := hooks.Script(args[0], string(hookSrc))
		if err != nil {
			return err
		}
		ctx.AddHook(args[0], hookObject)
		return nil
	}
	return errors.New("invalid syntax")
}
