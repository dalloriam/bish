package lib

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

func init() {
	registry = append(registry, getShellCallFn)
}

func getShellCallFn(sh CmdExecFn) *starlark.Builtin {
	return starlark.NewBuiltin("sh", func(t *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid syntax")
		}
		rawLine, ok := args[0].(starlark.String)
		if !ok {
			return nil, errors.New("shell call argument is not string")
		}

		output, succeeded, err := exec(string(rawLine), sh)
		if err != nil {
			return nil, err
		}

		if !succeeded {
			return starlark.String(fmt.Sprintf("error: %s", output)), nil
		}

		return starlark.String(output), nil
	})
}
