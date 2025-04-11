package system

import (
	"fmt"
	"os"
	"os/exec"
)

func Format(path, filesystem string) error {
	lookup, err := FindPath("mkfs." + filesystem)
	if err != nil {
		return fmt.Errorf("formatting tool for '%s' not found: %w", filesystem, err)
	}

	cmd := exec.Command(lookup, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format '%s' with filesystem '%s': %w", path, filesystem, err)
	}

	return nil
}
