package bish

import "io"

// ShellBackend abstracts the shell output interface for use by the shell engine.
type ShellBackend interface {
	UpdatePrompt(string)
	ReadLine() (string, error)
	Stdout() io.Writer
	Stderr() io.Writer
	Stdin() io.Reader
	Close() error
}
