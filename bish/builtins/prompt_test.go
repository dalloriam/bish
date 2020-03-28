package builtins_test

import (
	"github.com/dalloriam/bish/bish/constants"
	"github.com/dalloriam/bish/bish/state"
)

func validatePrompt(ctx *state.State, args []string) bool {
	if len(args) != 2 {
		return false
	}

	prompt, ok := ctx.GetKey(constants.PromptKey, constants.PromptKey)
	if !ok {
		return false
	}

	return prompt == args[1]
}
