package command

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/dalloriam/bish/bish/state"

	"github.com/dalloriam/bish/bish/builtins"
)

// Executable represents an actually-executable command.
type Executable struct {
	Args []Argument
	Ctx  *state.State

	Shell  bool
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader

	cmd *exec.Cmd
	buf bytes.Buffer
}

// Bind binds the provided inputs & outputs to the current command.
func (c *Executable) Bind(stdin io.Reader, stdout, stderr io.Writer) {
	c.StdOut = stdout
	c.StdErr = stderr
	c.StdIn = stdin
}

func (c *Executable) parseArguments() ([]string, error) {
	var args []string
	for _, a := range c.Args {
		aVal, err := a.Evaluate()
		if err != nil {
			return nil, err
		}
		aVals, err := ProcessArg(aVal, c.Ctx)
		if err != nil {
			return nil, err
		}
		args = append(args, aVals...)
	}
	return args, nil
}

// Start starts the command asynchronously.
func (c *Executable) Start() error {
	// Render the command arguments
	args, err := c.parseArguments()
	if err != nil {
		return err
	}
	_, ok, err := c.tryForBuiltIn(args)
	if ok {
		return err
	}
	// Run hooks.
	for _, hk := range c.Ctx.Hooks() {
		args = hk.Apply(args)
	}
	c.cmd = exec.Command(args[0], args[1:]...)

	var stdout io.Writer
	if c.Shell {
		stdout = c.StdOut
	} else {
		stdout = &c.buf
	}

	// Set the correct output device.
	c.cmd.Stderr = c.StdErr
	c.cmd.Stdout = stdout
	c.cmd.Stdin = c.StdIn

	return c.cmd.Start()
}

// Wait awaits the command.
func (c *Executable) Wait() (string, error) {
	if c.cmd == nil {
		data, err := ioutil.ReadAll(&c.buf)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	if c.Shell {
		return "", c.cmd.Wait()
	}
	data, err := ioutil.ReadAll(&c.buf)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Executable) nativeExec(args []string) (string, error) {
	// Run hooks.
	for _, hk := range c.Ctx.Hooks() {
		args = hk.Apply(args)
	}

	if c.Shell {
		if err := c.Start(); err != nil {
			return "", err
		}
		return c.Wait()
	}

	cmd := exec.Command(args[0], args[1:]...)
	bBuf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bBuf)), nil
}

func (c *Executable) tryForBuiltIn(args []string) (string, bool, error) {
	switch args[0] {
	case builtins.CdName:
		return "", true, builtins.ChangeDirectory(args[1:])
	case builtins.ExitName:
		return "", true, builtins.Exit()
	case builtins.AliasName:
		return "", true, builtins.Alias(c.Ctx, args[1:])
	case builtins.PromptName:
		return "", true, builtins.Prompt(c.Ctx, args[1])
	case builtins.HookName:
		return "", true, builtins.Hook(c.Ctx, args[1:])
	case builtins.SetEnvName:
		return "", true, builtins.SetEnv(args[1:])
	}
	return "", false, nil
}

// Signal sends an OS signal to the currently running command.
func (c *Executable) Signal(s os.Signal) {
	if c.cmd != nil && c.cmd.ProcessState != nil && !c.cmd.ProcessState.Exited() {
		if err := c.cmd.Process.Signal(s); err != nil {
			panic(err)
		}
	}
}

// Evaluate runs the command and waits for it to complete.
func (c *Executable) Evaluate() (string, error) {
	args, err := c.parseArguments()
	if err != nil {
		return "", err
	}

	out, ok, err := c.tryForBuiltIn(args)
	if ok {
		return out, err
	}

	return c.nativeExec(args)
}
