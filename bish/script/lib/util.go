package lib

import (
	"bytes"
	"io"
	"io/ioutil"
)

// A CmdExecFn is a callback to execute commands from the skylark API.
type CmdExecFn func(string, io.Reader, io.Writer, io.Writer) error

func exec(cmd string, sh CmdExecFn) (string, bool, error) {
	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := sh(string(cmd), &stdin, &stdout, &stderr); err != nil {
		all, err := ioutil.ReadAll(&stderr)
		if err != nil {
			return "", false, err
		}
		return string(all), false, nil
	}

	all, err := ioutil.ReadAll(&stdout)
	if err != nil {
		return "", false, err
	}

	return string(all), true, nil
}
