package hooks

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

var t = &starlark.Thread{
	Name:  "hook",
	Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
}

type scriptHook struct {
	Name   string
	Source string

	matchCall starlark.Value
	applyCall starlark.Value
}

// Script returns a new script hook from provided source code.
func Script(name string, source string) (Hook, error) {
	s := &scriptHook{
		Name:   name,
		Source: source,
	}

	globals, err := starlark.ExecFile(t, fmt.Sprintf("%s.star", name), source, nil)
	if err != nil {
		return nil, err
	}

	if globals.Has("apply") {
		s.applyCall = globals["apply"]
	} else {
		return nil, errors.New("no apply() function defined")
	}

	return s, nil
}

func (j *scriptHook) runCall(call starlark.Value, args []string) (starlark.Value, error) {
	var argValues []starlark.Value
	for _, arg := range args {
		argValues = append(argValues, starlark.String(arg))
	}
	v, err := starlark.Call(t, call, starlark.Tuple{starlark.NewList(argValues)}, nil)
	return v, err
}

// Apply applies the script hook.
func (j *scriptHook) Apply(args []string) []string {
	if j.applyCall == nil {
		return args
	}
	v, err := j.runCall(j.applyCall, args)
	if err != nil {
		return args
	}

	lst, ok := v.(*starlark.List)
	if !ok {
		return args
	}

	var newArgs []string
	for i := 0; i < lst.Len(); i++ {
		newArgs = append(newArgs, string(lst.Index(i).(starlark.String)))
	}
	return newArgs
}
