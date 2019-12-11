package hooks

type Hook interface {
	Apply(args []string) []string
}
