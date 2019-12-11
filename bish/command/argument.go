package command

import (
	"fmt"

	"go.starlark.net/starlark"
)

var t = &starlark.Thread{
	Name:  "hook",
	Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
}

type Argument interface {
	Evaluate() (string, error)
}

type StringArgument struct {
	Value string
}

func (a *StringArgument) Evaluate() (string, error) { return a.Value, nil }

func (a *StringArgument) String() string { return a.Value }

type PythonArgument struct {
	Source string
}

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
