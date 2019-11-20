package command

import (
	"bytes"
	"fmt"
	"github.com/dalloriam/bish/bish/builtins"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

type CommandTree struct {
	Args []Argument

	Shell  bool
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader

	cmd *exec.Cmd
	buf bytes.Buffer
}

func (c *CommandTree) Bind(stdin io.Reader, stdout, stderr io.Writer) {
	c.StdOut = stdout
	c.StdErr = stderr
	c.StdIn = stdin
}

func (c *CommandTree) parseArguments() ([]string, error) {
	var args []string
	for _, a := range c.Args {
		aVal, err := a.Evaluate()
		if err != nil {
			return nil, err
		}
		args = append(args, aVal)
	}
	return args, nil
}

func (c *CommandTree) Start() error {
	// Render the command arguments
	args, err := c.parseArguments()
	if err != nil {
		return err
	}
	fmt.Println("Args: ", args)
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

func (c *CommandTree) Wait() (string, error) {
	if c.Shell {
		return "", c.cmd.Wait()
	}
	data, err := ioutil.ReadAll(&c.buf)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *CommandTree) nativeExec(args []string) (string, error) {
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

func (c *CommandTree) tryForBuiltIn(args []string) (string, bool, error) {
	switch args[0] {
	case builtins.CdName:
		return "", true, builtins.ChangeDirectory(args[1])
	case builtins.ExitName:
		return "", true, builtins.Exit()
	}
	return "", false, nil
}

func (c *CommandTree) Evaluate() (string, error) {
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