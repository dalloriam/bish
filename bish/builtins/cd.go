package builtins

import (
	"os"
)

const CdName = "cd"

// ChangeDirectory changes the directory.
func ChangeDirectory(newDirectory string) error {
	return os.Chdir(newDirectory)
}
