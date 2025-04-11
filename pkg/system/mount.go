package system

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

func IsMounted(path string) (bool, error) {
	mounts, err := os.Open("/proc/mounts")
	if err != nil {
		return false, fmt.Errorf("failed to access '/proc/mounts': %w", err)
	}
	defer mounts.Close()

	content, err := io.ReadAll(mounts)
	if err != nil {
		return false, fmt.Errorf("failed to read '/proc/mounts': %w", err)
	}

	return strings.Contains(string(content), " "+path+" "), nil
}

func MountImage(path, volume, filesystem string) error {
	var flags uintptr = syscall.MS_NOATIME

	// TODO: Additional mount options
	data := ""

	if err := syscall.Mount(path, volume, filesystem, flags, data); err != nil {
		return fmt.Errorf("failed to mount '%s': %w", path, err)
	}

	return nil
}

func UmountImage(path string) error {
	if err := syscall.Unmount(path, 0); err != nil {
		if err := syscall.Unmount(path, syscall.MNT_FORCE); err != nil {
			return fmt.Errorf("failed to umount '%s': %w", path, err)
		}
	}

	return nil
}
