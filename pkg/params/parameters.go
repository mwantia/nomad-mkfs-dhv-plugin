package params

type VolumeParameters struct {
	Filesystem   string `json:"filesystem"`
	BlockSize    string `json:"block_size"`
	MountOptions string `json:"mount_options"`
	ReadOnly     bool   `json:"read_only"`
}

func NewDefault() *VolumeParameters {
	return &VolumeParameters{
		Filesystem:   "ext4",
		BlockSize:    "4k",
		MountOptions: "noatime,nodiratime",
		ReadOnly:     false,
	}
}
