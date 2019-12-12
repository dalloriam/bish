package builtins

import "github.com/dalloriam/bish/bish/state"

// Name of the prompt builtin.
const (
	PromptName = "prompt"
)

// Prompt sets the user prompt in the context.
func Prompt(ctx *state.State, val string) error {
	ctx.SetKey("prompt", "prompt", val)
	return nil
}
