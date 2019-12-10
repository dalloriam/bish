package completion

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func dirExist(path string) bool {
	inf, err := os.Stat(path)
	if err != nil {
		return false
	}

	return inf.IsDir()
}

func listPathOptions(currentArgument string) []string {
	var currentArgDir string
	if dirExist(currentArgument) {
		currentArgDir = filepath.Clean(currentArgument)
	} else {
		currentArgDir = filepath.Dir(filepath.Clean(currentArgument))
	}

	fileInfos, err := ioutil.ReadDir(currentArgDir)
	if err != nil {
		// TODO: Log somewhere.
		return nil
	}

	includeHidden := strings.HasPrefix(filepath.Base(currentArgument), ".")

	var options []string
	for _, fileInfo := range fileInfos {
		if strings.HasPrefix(fileInfo.Name(), ".") && !includeHidden {
			continue
		}
		value := filepath.Join(currentArgDir, fileInfo.Name())
		if fileInfo.IsDir() {
			value = value + "/"
		}
		options = append(options, value)
	}

	return options
}
