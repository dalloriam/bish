package script

import (
	"errors"

	"go.starlark.net/starlark"
)

// A Hook calls a starlark function before command execution.
type Hook struct {
	Name    string
	Source  string
	runtime *Runtime

	applyCall starlark.Value
}

// NewHook returns a new script hook from provided source code.
func NewHook(name string, source string, runtime *Runtime) (*Hook, error) {
	s := &Hook{
		Name:    name,
		Source:  source,
		runtime: runtime,
	}

	globals, err := s.runtime.ExecSandbox(s.Name, s.Source)
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

func (j *Hook) runCall(call starlark.Value, args []string) (starlark.Value, error) {
	var argValues []starlark.Value
	for _, arg := range args {
		argValues = append(argValues, starlark.String(arg))
	}
	return j.runtime.CallRaw(j.applyCall, []starlark.Value{starlark.NewList(argValues)}, nil)
}

// Apply applies the script hook.
func (j *Hook) Apply(args []string) []string {
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
