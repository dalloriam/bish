package lib

import "go.starlark.net/starlark"

type libFn func(fn CmdExecFn) *starlark.Builtin

var registry []libFn

// GetBuiltins generates the standard library builtins for skylark scripts.
func GetBuiltins(fn CmdExecFn) []*starlark.Builtin {
	var builtins []*starlark.Builtin

	for _, v := range registry {
		builtins = append(builtins, v(fn))
	}
	return builtins
}
