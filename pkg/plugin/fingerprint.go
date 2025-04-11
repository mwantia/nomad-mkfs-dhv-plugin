package plugin

import (
	"encoding/json"
	"fmt"
)

func Fingerprint() error {
	resp := FingerprintResponse{
		Version: Version,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	fmt.Print(string(json))
	return nil
}
