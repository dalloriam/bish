package bish

type ShellBackend interface {
	ReadLine() (string, error)
	Stdout(a string)
	Stderr(a string)
	Close() error
}
