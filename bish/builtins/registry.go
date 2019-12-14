package builtins

import "github.com/dalloriam/bish/bish/state"

// A Builtin function is shell built-in behavior that is executed like a command.
type Builtin func(ctx *state.State, args []string) (string, error)

var registry map[string]Builtin = make(map[string]Builtin)

// TryBuiltIns attempts to execute known builtins for the current command, returning
// `true` as second return value if one was found.
func TryBuiltIns(ctx *state.State, args []string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, nil
	}

	cmd := args[0]
	args = args[1:]
	if bt, ok := registry[cmd]; ok {
		o, err := bt(ctx, args)
		return o, true, err
	}
	return "", false, nil
}
