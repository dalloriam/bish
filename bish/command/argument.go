package command

import "errors"

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
	return "", errors.New("pythonArg not implemented")
}
