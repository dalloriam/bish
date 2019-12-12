package main

import (
	"io"
	"os"

	"github.com/chzyer/readline"
)

// TerminalBackend implements the backend interface expected by the shell core.
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

// UpdatePrompt updates the user prompt.
func (t *TerminalBackend) UpdatePrompt(prompt string) {
	t.rl.SetPrompt(prompt)
}

// SetConfig sets the autocomplete engine.
func (t *TerminalBackend) SetConfig(c readline.AutoCompleter) {
	t.rl.Config.AutoComplete = c
}

// Stderr returns the preferred stderr stream.
func (t *TerminalBackend) Stderr() io.Writer {
	return os.Stderr
}

// Stdin returns the preferred stdin stream.
func (t *TerminalBackend) Stdin() io.Reader {
	return os.Stdin
}

// Stdout returns the preferred stdout stream.
func (t *TerminalBackend) Stdout() io.Writer {
	return os.Stdout
}

// ReadLine reads a single line from the user.
func (t *TerminalBackend) ReadLine() (string, error) {
	return t.rl.Readline()
}

// Close terminates the shell backend.
func (t *TerminalBackend) Close() error {
	return t.rl.Close()
}
