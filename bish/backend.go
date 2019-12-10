package bish

import "io"

type ShellBackend interface {
	UpdatePrompt(string)
	ReadLine() (string, error)
	Stdout() io.Writer
	Stderr() io.Writer
	Stdin() io.Reader
	Close() error
}
