package plugin

const Version = "1.0.0"

type FingerprintResponse struct {
	Version string `json:"version"`
}

type VolumeCreateResponse struct {
	Path  string `json:"path"`
	Bytes int64  `json:"bytes"`
}
