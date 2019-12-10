package builtins

const (
	PromptName = "prompt"
)

// Prompt sets the user prompt in the context.
func Prompt(ctx ShellContext, val string) error {
	ctx.SetKey("prompt", "prompt", val)
	return nil
}
