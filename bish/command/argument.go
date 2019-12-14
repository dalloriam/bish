package command

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
