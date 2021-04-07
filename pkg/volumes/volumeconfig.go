package volumes

type VolumeConfig struct {
	// Name of the volume; current only tested on macOS.
	Name *string `yaml:"name"`
	// Mount path of the volume; tested on Linux and macOS,
	// much more precise than Name.
	MountPath *string `yaml:"mount_path"`
	// Actual device name of the volume,
	// ex. "/dev/sda1" (Linux) or "/dev/disk5s1" (macOS)
	Device *string `yaml:"device"`
}
