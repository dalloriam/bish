package completion

import (
	"io/ioutil"
	"os"
	"strings"
)

func listExecutables() []string {
	var out []string
	// TODO: Usage statistics.

	pathEnv := os.Getenv("PATH")
	for _, dirPath := range strings.Split(pathEnv, ":") {
		listing, err := ioutil.ReadDir(dirPath)
		if err != nil {
			// TODO: Log somewhere.
			continue
		}

		for _, fileName := range listing {
			if (fileName.Mode() & 0111) != 0 {
				out = append(out, fileName.Name())
			}
		}
	}

	return out
}
