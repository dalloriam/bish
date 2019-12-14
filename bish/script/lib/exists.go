package lib

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

func init() {
	registry = append(registry, exists)
}

func exists(sh CmdExecFn) *starlark.Builtin {
	return starlark.NewBuiltin("exists", func(t *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid syntax")
		}
		rawLine, ok := args[0].(starlark.String)
		if !ok {
			return nil, errors.New("shell call argument is not string")
		}
		_, succeeded, err := exec(fmt.Sprintf("which %s ", string(rawLine)), sh)
		if err != nil {
			return nil, err
		}

		return starlark.Bool(succeeded), nil
	})
}
