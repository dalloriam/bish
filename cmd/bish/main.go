package main

import (
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/dalloriam/bish/bish"
)

type TerminalBackend struct {
	rl *readline.Instance
}

func newBackend() (*TerminalBackend, error) {
	rl, err := readline.New("➤ ")
	if err != nil {
		return nil, err
	}

	return &TerminalBackend{rl: rl}, nil
}

func (t *TerminalBackend) Stderr() io.Writer {
	return os.Stderr
}

func (t *TerminalBackend) Stdin() io.Reader {
	return os.Stdin
}

func (t *TerminalBackend) Stdout() io.Writer {
	return os.Stdout
}

func (t *TerminalBackend) ReadLine() (string, error) {
	return t.rl.Readline()
}

func (t *TerminalBackend) Close() error {
	return t.rl.Close()
}

func shellStart() {
	backend, err := newBackend()
	if err != nil {
		panic(err)
	}
	defer backend.Close()

	shell := bish.New(backend)
	shell.Start()
}

func main() {
	shellStart()
}
