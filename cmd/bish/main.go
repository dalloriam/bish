package main

import (
	"os"

	"github.com/dalloriam/bish/bish/config"

	"github.com/dalloriam/bish/bish"
)

func main() {
	os.Setenv("TERM", "xterm-kitty")
	shell := bish.New(config.IOConfig{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	shell.Start()
}
