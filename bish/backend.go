package bish

import "io"

type ShellBackend interface {
	ReadLine() (string, error)
	Stdout() io.Writer
	Stderr() io.Writer
	Stdin() io.Reader
	Close() error
}
