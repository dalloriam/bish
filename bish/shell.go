package bish

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/dalloriam/bish/bish/state"

	"github.com/dalloriam/bish/bish/completion"

	"github.com/dalloriam/bish/bish/command"
	. "github.com/logrusorgru/aurora"
)

// Shell is the root bish shell struct.
type Shell struct {
	backend ShellBackend
	ctx     *state.State

	CompletionProvider *completion.Completer

	promptExpr string
}

// New initializes & returns a new shell instance with the provided IO writers.
func New(backend ShellBackend) *Shell {
	ctx := state.New()
	prompt := "{{.Username}}@{{.Hostname}} [{{.Cwd}}]\n➤ "
	return &Shell{backend: backend, ctx: ctx, CompletionProvider: completion.New(ctx), promptExpr: prompt}
}

// prompt returns a slice of prompt lines.
func (s *Shell) prompt() []string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	relCwd, err := filepath.Rel(u.HomeDir, cwd)
	if err != nil {
		relCwd = cwd
	} else {
		relCwd = filepath.Join("~", relCwd)
	}

	promptTemplate := "{{.Username}}@{{.Hostname}} ➤ "
	promptTemplateInt, ok := s.ctx.GetKey("prompt", "prompt")
	if ok {
		promptTemplate = promptTemplateInt.(string)
	}

	pCtx := promptContext{Username: u.Username, Cwd: relCwd, Hostname: "phoenix"}
	v, err := pCtx.Render(promptTemplate)
	if err != nil {
		// TODO: Handle invalid prompts
		panic(err)
	}
	return v
}

func (s *Shell) err(err error) {
	_, err = s.backend.Stderr().Write([]byte(fmt.Sprintf("! %s\n", Red(err.Error()))))
	return
}

func (s *Shell) stdout(d string) (err error) {
	_, err = s.backend.Stdout().Write([]byte(d))
	return
}

// Start starts the main shell loop.
func (s *Shell) Start() {
	for {
		promptList := s.prompt()
		if err := s.stdout(strings.Join(promptList[:len(promptList)], "\n")); err != nil {
			s.err(err)
			continue
		}
		s.backend.UpdatePrompt(promptList[len(promptList)-1])

		rawLine, err := s.backend.ReadLine()
		if err != nil {
			s.err(err)
			continue
		}
		req := command.CommandRequest{
			Context:   s.ctx,
			UserInput: rawLine,
			Stdin:     s.backend.Stdin(),
			Stdout:    s.backend.Stdout(),
			Stderr:    s.backend.Stderr(),
		}

		if err := req.Execute(); err != nil {
			s.err(err)
			continue
		}
	}
}
