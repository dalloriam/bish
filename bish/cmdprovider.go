package bish

// A CommandReader reads a command.
type CommandReader interface {
	Read() (string, error)
}
