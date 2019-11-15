package bish

type ShellBackend interface {
	ReadLine() (string, error)
	Stdout(a ...interface{})
	Stderr(a ...interface{})
	Close() error
}
