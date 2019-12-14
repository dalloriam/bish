package bish

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/dalloriam/bish/bish/script"

	"github.com/dalloriam/bish/bish/state"

	"github.com/dalloriam/bish/bish/completion"

	"github.com/dalloriam/bish/bish/command"
	cl "github.com/logrusorgru/aurora"
)

// Shell is the root bish shell struct.
type Shell struct {
	backend ShellBackend

	ctx     *state.State
	runtime *script.Runtime

	CompletionProvider *completion.Completer

	promptExpr string
}

// New initializes & returns a new shell instance with the provided IO writers.
func New(backend ShellBackend) *Shell {
	prompt := "{{.Username}}@{{.Hostname}} [{{.Cwd}}]\n➤ "
	sh := &Shell{backend: backend, promptExpr: prompt}
	ctx := state.New(sh.runCmd)
	sh.ctx = ctx
	sh.CompletionProvider = completion.New(ctx)
	return sh
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

	hostname, err := os.Hostname()
	if err != nil {
		// TODO: Handle gracefully.
		panic(err)
	}
	splittedHostname := strings.Split(hostname, ".")
	if len(splittedHostname) > 0 {
		hostname = splittedHostname[0]
	}

	promptTemplate := "{{.Username}}@{{.Hostname}} ➤ "
	promptTemplateInt, ok := s.ctx.GetKey("prompt", "prompt")
	if ok {
		promptTemplate = promptTemplateInt.(string)
	}

	pCtx := promptContext{Username: u.Username, Cwd: relCwd, Hostname: hostname}
	v, err := pCtx.Render(promptTemplate)
	if err != nil {
		// TODO: Handle invalid prompts
		panic(err)
	}
	return v
}

func (s *Shell) warn(str string) {
	s.backend.Stdout().Write([]byte(fmt.Sprintf("! %s\n", cl.Yellow(str))))
	return
}

func (s *Shell) err(err error) {
	_, err = s.backend.Stderr().Write([]byte(fmt.Sprintf("! %s\n", cl.Red(err.Error()))))
	return
}

func (s *Shell) stdout(d string) (err error) {
	_, err = s.backend.Stdout().Write([]byte(d))
	return
}

func (s *Shell) execScript(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = s.ctx.ScriptRuntime().Exec(path, string(data))
	return err
}

func (s *Shell) importBishFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for _, l := range strings.Split(string(data), "\n") {
		if err := s.runCmd(l, s.backend.Stdin(), s.backend.Stdout(), s.backend.Stderr()); err != nil {
			return err
		}
	}

	return nil
}

func (s *Shell) runCmd(rawLine string, stdin io.Reader, stdout, stderr io.Writer) error {
	req := command.Request{
		Context:   s.ctx,
		UserInput: rawLine,
		Stdin:     stdin,
		Stdout:    stdout,
		Stderr:    stderr,
	}
	return req.Execute()
}

func (s *Shell) configure() error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	if err := os.Chdir(configDir); err != nil {
		return err
	}
	defer os.Chdir(curDir)

	cfgFiles, err := getConfigFiles()
	if err != nil {
		return err
	}

	for _, cfgFile := range cfgFiles {
		if strings.HasSuffix(cfgFile, "star") {
			if err := s.execScript(cfgFile); err != nil {
				return err
			}
		} else if strings.HasSuffix(cfgFile, "bish") {
			if err := s.importBishFile(cfgFile); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unimportable: %s", cfgFile)
		}
	}
	return nil
}

// Start starts the main shell loop.
func (s *Shell) Start() {

	if err := s.configure(); err != nil {
		s.err(err)
	}

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

		if err := s.runCmd(rawLine, s.backend.Stdin(), s.backend.Stdout(), s.backend.Stderr()); err != nil {
			s.err(err)
			continue
		}
	}
}
