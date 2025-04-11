package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mwantia/nomad-mkfs-host-volume-plugin/pkg/params"
)

func Delete() error {
	volumesDir := os.Getenv("DHV_VOLUMES_DIR")
	volumeID := os.Getenv("DHV_VOLUME_ID")
	createdPath := os.Getenv("DHV_CREATED_PATH")

	path := createdPath
	if path == "" {
		if volumesDir == "" {
			return fmt.Errorf("DHV_VOLUMES_DIR must not be empty")
		}
		if volumeID == "" {
			return fmt.Errorf("DHV_VOLUME_ID must not be empty")
		}

		path = filepath.Join(volumesDir, volumeID)
	}

	log.Printf("Deleting volume at %s", path)

	parameters := params.NewDefault()
	paramsJson := os.Getenv("DHV_PARAMETERS")

	if paramsJson != "" {
		if err := json.Unmarshal([]byte(paramsJson), &parameters); err != nil {
			log.Printf("Warning: Unable to parse parameters: %v, using defaults", err)
		}
	}

	if isMounted(path) {
		log.Printf("Unmounting %s", path)
		umountCmd := exec.Command("/usr/bin/umount", path)

		if err := umountCmd.Run(); err != nil {
			return fmt.Errorf("failed to unmount: %v", err)
		}
	}

	if err := os.RemoveAll(path); err != nil {
		log.Printf("Warning: Failed to remove directory: %v", err)
	}
	if err := os.Remove(fmt.Sprintf("%s.%s", path, parameters.Filesystem)); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Failed to remove ext4 image: %v", err)
	}

	return nil
}
