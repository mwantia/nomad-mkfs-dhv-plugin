package plugin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mwantia/nomad-mkfs-dhv-plugin/pkg/config"
	"github.com/mwantia/nomad-mkfs-dhv-plugin/pkg/system"
)

func Delete(cfg config.DynamicHostVolumeConfig) error {
	path := cfg.CreatedPath
	if path == "" {
		if cfg.VolumesDir == "" {
			return fmt.Errorf("variable 'DHV_VOLUMES_DIR' must not be empty when 'DHV_CREATED_PATH' is not provided")
		}
		if cfg.VolumeID == "" {
			return fmt.Errorf("variable 'DHV_VOLUME_ID' must not be empty when 'DHV_CREATED_PATH' is not provided")
		}

		path = filepath.Join(cfg.VolumesDir, cfg.VolumeID)
	}

	log.Printf("Deleting volume at %s", path)

	mounted, err := system.IsMounted(path)
	if err != nil {
		return fmt.Errorf("failed to check mount for '%s': %w", path, err)
	}

	if mounted {
		log.Printf("Unmounting %s...", path)
		if err := system.UmountImage(path); err != nil {
			log.Printf("Warning: Unable to unmount '%s': %v", path, err)
		}
	}

	if err := os.RemoveAll(path); err != nil {
		log.Printf("Warning: Failed to remove '%s': %v", path, err)
	}

	if err := os.Remove(path + ".img"); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Failed to remove '%s': %v", path+".img", err)
	}

	return nil
}
