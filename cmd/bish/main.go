package main

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/dalloriam/bish/bish"
	. "github.com/logrusorgru/aurora"
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

func (t *TerminalBackend) Stderr(a string) {
	fmt.Fprint(t.rl.Stderr(), Red(a))
}

func (t *TerminalBackend) Stdout(a string) {
	fmt.Fprint(t.rl.Stdout(), a)
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
