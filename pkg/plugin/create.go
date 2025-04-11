package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mwantia/nomad-mkfs-host-volume-plugin/pkg/params"
)

func Create() error {
	volumesDir := os.Getenv("DHV_VOLUMES_DIR")
	volumeID := os.Getenv("DHV_VOLUME_ID")
	capacityMinBytesStr := os.Getenv("DHV_CAPACITY_MIN_BYTES")

	if volumesDir == "" {
		return fmt.Errorf("variable 'DHV_VOLUMES_DIR' must not be empty")
	}
	if volumeID == "" {
		return fmt.Errorf("variable 'DHV_VOLUME_ID' must not be empty")
	}

	parameters := params.NewDefault()
	paramsJson := os.Getenv("DHV_PARAMETERS")

	if paramsJson != "" {
		if err := json.Unmarshal([]byte(paramsJson), &parameters); err != nil {
			log.Printf("Warning: Unable to parse parameters: %v, using defaults", err)
		}
	}

	volumePath := filepath.Join(volumesDir, volumeID)
	imagePath := fmt.Sprintf("%s.%s", volumePath, parameters.Filesystem)

	if err := os.MkdirAll(volumePath, 0o755); err != nil {
		return fmt.Errorf("failed to create volume directory: %v", err)
	}

	capacityMinBytes, err := strconv.ParseInt(capacityMinBytesStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse capacity: %v", err)
	}
	if capacityMinBytes <= 0 {
		return fmt.Errorf("minimum capacity must be greater than zero")
	}

	capacityMB := capacityMinBytes / (1024 * 1024)
	if capacityMB <= 0 {
		capacityMB = 1 // Ensure at least 1MB
	}

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("Creating filesystem image at %s", imagePath)

		ddCmd := exec.Command("/usr/bin/dd", "if=/dev/zero", "of="+imagePath, "bs=1M", "count="+strconv.FormatInt(capacityMB, 10))
		ddCmd.Stderr = os.Stderr
		if err := ddCmd.Run(); err != nil {
			return fmt.Errorf("failed to create filesystem image: %v", err)
		}

		mkfsCmd := exec.Command("/usr/sbin/mkfs."+parameters.Filesystem, imagePath)
		mkfsCmd.Stderr = os.Stderr

		if err := mkfsCmd.Run(); err != nil {
			return fmt.Errorf("failed to format filesystem: %v", err)
		}
	} else {
		log.Printf("Using existing filesystem image at %s", imagePath)
	}

	if !isMounted(volumePath) {
		log.Printf("Mounting filesystem at %s", volumePath)
		mountCmd := exec.Command("/usr/bin/mount", imagePath, volumePath)

		if err := mountCmd.Run(); err != nil {
			return fmt.Errorf("failed to mount filesystem: %v", err)
		}
	} else {
		log.Printf("Filesystem already mounted at %s", volumePath)
	}

	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		return fmt.Errorf("failed to get filesystem size: %v", err)
	}
	actualBytes := fileInfo.Size()

	response := VolumeCreateResponse{
		Path:  volumePath,
		Bytes: actualBytes,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	fmt.Print(string(jsonResponse))
	return nil
}

func isMounted(path string) bool {
	cmd := exec.Command("/usr/bin/mount")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), " "+path+" ")
}
