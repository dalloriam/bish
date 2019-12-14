package command

import (
	"github.com/dalloriam/bish/bish/state"
)

// PythonArgument represents an embedded python subcommand.
type PythonArgument struct {
	Source string
	State  *state.State
}

// Evaluate evaluates the python subcommand.
func (a *PythonArgument) Evaluate() (string, error) {
	return a.State.ScriptRuntime().Eval(a.Source)
}
