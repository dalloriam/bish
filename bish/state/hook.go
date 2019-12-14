package state

// A Hook pre-processes command arguments.
type Hook interface {
	Apply(args []string) []string
}
