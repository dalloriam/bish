package bish

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

var definedConfigDirectories = []string{"conf.d", "functions"}

func isImportable(f os.FileInfo) bool {
	return !f.IsDir() && (strings.HasSuffix(f.Name(), ".star") || strings.HasSuffix(f.Name(), ".bish"))
}

func getConfigDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	configDir := path.Join(user.HomeDir, ".config", "bish")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}
	return configDir, nil
}

func getConfigFiles() ([]string, error) {

	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	var files []string
	for _, definedDirectory := range definedConfigDirectories {
		dirPath := path.Join(configDir, definedDirectory)

		listing, err := ioutil.ReadDir(dirPath)
		if err != nil {
			continue
		}

		for _, f := range listing {
			if isImportable(f) {
				files = append(files, path.Join(dirPath, f.Name()))
			}
		}
	}

	return files, nil
}
