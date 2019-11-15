package config

import "io"

// IOConfig holds the I/O specification of the shell instance.
type IOConfig struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
