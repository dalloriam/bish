package command

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// A PipeCommand represents a command piped into another.
type PipeCommand struct {
	SrcCommand Command
	DstCommand Command

	Shell  bool
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader

	buf       bytes.Buffer
	pipeWrite *io.PipeWriter
}

// Signal sends the OS signal to both commands in the pipe sequence.
func (c *PipeCommand) Signal(s os.Signal) {
	c.SrcCommand.Signal(s)
	c.DstCommand.Signal(s)
}

// Bind binds the provided stdin to the first command, and the stdout/stderr to the second command in the pipe sequence.
func (c *PipeCommand) Bind(stdin io.Reader, stdout, stderr io.Writer) {
	c.StdIn = stdin
	c.StdOut = stdout
	c.StdErr = stderr
}

// Start starts the pipe sequence asynchronously.
func (c *PipeCommand) Start() error {
	r, w := io.Pipe()
	c.pipeWrite = w

	var dstStdout io.Writer

	if c.Shell {
		dstStdout = c.StdOut
	} else {
		dstStdout = &c.buf
	}

	c.SrcCommand.Bind(c.StdIn, w, c.StdErr)
	c.DstCommand.Bind(r, dstStdout, c.StdErr)

	if err := c.SrcCommand.Start(); err != nil {
		return err
	}
	return c.DstCommand.Start()
}

// Wait awaits the pipe sequence.
func (c *PipeCommand) Wait() (string, error) {
	if _, err := c.SrcCommand.Wait(); err != nil {
		return "", err
	}

	if err := c.pipeWrite.Close(); err != nil {
		return "", err
	}

	if _, err := c.DstCommand.Wait(); err != nil {
		return "", err
	}

	var out string
	if !c.Shell {
		data, err := ioutil.ReadAll(&c.buf)
		if err != nil {
			return "", err
		}
		out = string(data)
	}
	return out, nil
}

// Evaluate runs & awaits the pipe sequence.
func (c *PipeCommand) Evaluate() (string, error) {
	if err := c.Start(); err != nil {
		return "", err
	}
	return c.Wait()
}
