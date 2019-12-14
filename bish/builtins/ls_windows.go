package builtins

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/dalloriam/bish/bish/state"
)

const (
	lsName  = "ls"
	dirName = "dir"
)

func init() {
	registry[lsName] = lsDir
	registry[dirName] = lsDir
}

// lsDir is a basic `ls` implementation for windows systems.
func lsDir(ctx *state.State, args []string) (string, error) {
	var dirPath string
	var err error
	if len(args) == 0 {
		dirPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	} else {
		dirPath = args[0]
	}
	listing, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	var directoryNames []string

	for _, v := range listing {
		name := v.Name()
		if strings.HasPrefix(name, ".") {
			// Skip hidden files.
			// TODO: Add flag for hidden.
			continue
		}

		if v.IsDir() {
			name += "/"
		}

		directoryNames = append(directoryNames, name)
	}

	sort.Strings(directoryNames)

	out := strings.Join(directoryNames, "\n")
	return out, nil
}
