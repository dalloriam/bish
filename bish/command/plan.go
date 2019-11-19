package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

/*
Command  :- Argument [Argument...] [>> file]
		 |- Command | Command

Argument :- String
		 |- ( Command )
		 |- < String >
*/

type Argument interface {
	Evaluate() (string, error)
}

type PipeCommand struct {
	SrcCommand *CommandTree
	DstCommand *CommandTree

	Shell  bool
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader
}

func (c *PipeCommand) Evaluate() (string, error) {

	r, w := io.Pipe()

	// Actual piping part.
	c.SrcCommand.StdOut = w
	c.DstCommand.StdIn = r

	// Set the correct output device.
	if c.StdErr != nil {
		c.SrcCommand.StdErr = c.StdErr
		c.DstCommand.StdErr = c.StdErr
	}

	if c.StdIn != nil {
		c.SrcCommand.StdIn = c.StdIn
	}

	var b2 bytes.Buffer
	if c.Shell {
		c.DstCommand.StdOut = c.StdOut
	} else {
		c.DstCommand.StdOut = &b2
	}

	c1, err := c.SrcCommand.Start()
	if err != nil {
		return "", err
	}
	c2, err := c.DstCommand.Start()
	if err != nil {
		return "", err
	}

	if err := c1.Wait(); err != nil {
		return "", err
	}

	w.Close()

	if err := c2.Wait(); err != nil {
		return "", err
	}

	if c.Shell {
		return "", nil
	}

	data, err := ioutil.ReadAll(&b2)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type CommandTree struct {
	Args []Argument

	Shell  bool
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader
}

func (c *CommandTree) ParseArguments() ([]string, error) {
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

func (c *CommandTree) Start() (*exec.Cmd, error) {
	// Render the command arguments
	args, err := c.ParseArguments()
	if err != nil {
		return nil, err
	}
	fmt.Println("Args: ", args)
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	if c.StdErr != nil {
		cmd.Stderr = c.StdErr
	}

	if c.StdOut != nil {
		cmd.Stdout = c.StdOut
	}

	if c.StdIn != nil {
		cmd.Stdin = c.StdIn
	}

	return cmd, cmd.Start()
}

func (c *CommandTree) Evaluate() (string, error) {
	if c.Shell {
		cmd, err := c.Start()
		if err != nil {
			return "", err
		}
		return "", cmd.Wait()
	}

	args, err := c.ParseArguments()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(args[0], args[1:]...)
	bBuf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bBuf)), nil
}

type StringArgument struct {
	Value string
}

func (a *StringArgument) Evaluate() (string, error) { return a.Value, nil }

type PythonArgument struct {
	Source string
}

func (a *PythonArgument) Evaluate() (string, error) {
	return "", errors.New("pythonArg not implemented")
}
