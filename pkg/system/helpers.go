package system

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func FindPath(file string) (string, error) {
	paths := []string{
		"/sbin",
		"/usr/sbin",
		"/bin",
		"/usr/bin",
		"/usr/local/bin",
		"/usr/local/sbin",
	}

	for _, path := range paths {
		full := filepath.Join(path, file)
		if lookup, err := exec.LookPath(full); err != nil {
			return lookup, nil
		}
	}

	return "", fmt.Errorf("file '%s' not found", file)
}
