package command

import (
	"bufio"
	"io"
)

// CommandReader streams command from an io.Reader.
type CommandReader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *CommandReader {
	return &CommandReader{reader: bufio.NewReader(r)}
}

func (r *CommandReader) Read() (*Command, error) {
	raw, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return ParseCommand(raw)
}
