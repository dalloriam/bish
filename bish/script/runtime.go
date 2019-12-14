package script

import (
	"fmt"

	"github.com/dalloriam/bish/bish/script/lib"
	"go.starlark.net/starlark"
)

// The Runtime is the global starlark runtime context.
type Runtime struct {
	state  starlark.StringDict
	execFn lib.CmdExecFn
}

// NewRuntime initializes the scripting runtime.
func NewRuntime(fn lib.CmdExecFn) *Runtime {
	r := &Runtime{
		state: make(map[string]starlark.Value),
	}

	for _, builtIn := range lib.GetBuiltins(fn) {
		r.state[builtIn.Name()] = builtIn
	}

	return r
}

// CallRaw calls a starlark function object.
func (r *Runtime) CallRaw(fn starlark.Value, args []starlark.Value, kwargs map[string]starlark.Value) (starlark.Value, error) {
	thread := getThread()
	defer putThread(thread)

	var tupleKwargs []starlark.Tuple
	for k, v := range kwargs {
		tupleKwargs = append(tupleKwargs, starlark.Tuple{
			starlark.String(k),
			v,
		})
	}

	returnVal, err := starlark.Call(
		thread,
		fn,
		starlark.Tuple(args),
		tupleKwargs,
	)

	return returnVal, err
}

// Call calls a function defined in-state by name.
func (r *Runtime) Call(fnName string, args []starlark.Value, kwargs map[string]starlark.Value) (starlark.Value, error) {
	fn, ok := r.state[fnName]
	if !ok {
		return nil, fmt.Errorf("unknown function: %s", fnName)
	}
	return r.CallRaw(fn, args, kwargs)
}

// EvalRaw evaluates a starlark expression without interpreting the output.
func (r *Runtime) EvalRaw(expr string) (starlark.Value, error) {
	thread := getThread()
	defer putThread(thread)

	return starlark.Eval(thread, "starlark_expression", expr, r.state)
}

// Eval evaluates a starlark expression as string.
func (r *Runtime) Eval(expr string) (string, error) {
	starlarkValue, err := r.EvalRaw(expr)
	if err != nil {
		return "", err
	}
	if s, ok := starlarkValue.(starlark.String); ok {
		return string(s), nil
	}
	return starlarkValue.String(), nil
}

// ExecSandbox runs the provided starlark script in an isolated state.
func (r *Runtime) ExecSandbox(name string, src string) (starlark.StringDict, error) {
	thread := getThread()
	defer putThread(thread)

	return starlark.ExecFile(thread, name, src, nil)
}

// Exec runs the provided starlark script in the global state.
func (r *Runtime) Exec(name string, src string) (starlark.StringDict, error) {
	thread := getThread()
	defer putThread(thread)

	globals, err := starlark.ExecFile(thread, name, src, r.state)
	if err != nil {
		return nil, err
	}

	// Merge exec'd globals with current state.
	for k, v := range globals {
		r.state[k] = v
	}

	return globals, nil
}

// Get fetches a value from the runtime state.
func (r *Runtime) Get(name string) (starlark.Value, bool) {
	v, ok := r.state[name]
	return v, ok
}
