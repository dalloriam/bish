package main

import (
	"github.com/dalloriam/bish/bish"
)

func main() {
	backend, err := newBackend()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := backend.Close(); err != nil {
			panic(err)
		}
	}()

	shell := bish.New(backend)
	backend.SetConfig(shell.CompletionProvider)
	shell.Start()
}
