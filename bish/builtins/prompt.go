package builtins

import (
	"errors"

	"github.com/dalloriam/bish/bish/constants"

	"github.com/dalloriam/bish/bish/state"
)

// Name of the prompt builtin.
const (
	promptName = "prompt"
)

func init() {
	registry[promptName] = prompt
}

func prompt(ctx *state.State, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid set prompt syntax")
	}
	ctx.SetKey(constants.PromptKey, constants.PromptKey, args[0])
	return "", nil
}
