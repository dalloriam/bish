package command

import (
	"fmt"

	"go.starlark.net/starlark"
)

// TODO: Clean.
var t = &starlark.Thread{
	Name:  "hook",
	Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
}

// Argument represents a shell argument that has a value.
type Argument interface {
	Evaluate() (string, error)
}

// StringArgument is a simple string arg.
type StringArgument struct {
	Value string
}

// Evaluate returns the string value of the arg.
func (a *StringArgument) Evaluate() (string, error) { return a.Value, nil }

// String returns the string value.
func (a *StringArgument) String() string { return a.Value }

// PythonArgument represents an embedded python subcommand.
type PythonArgument struct {
	Source string
}

// Evaluate evaluates the python subcommand.
func (a *PythonArgument) Evaluate() (string, error) {
	v, err := starlark.Eval(t, "arg", a.Source, nil)
	if err != nil {
		return "", err
	}
	if s, ok := v.(starlark.String); ok {
		return string(s), nil
	}
	return v.String(), nil
}
