package volumes

import (
	"context"
	"errors"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

func GetVolumePath(cfg VolumeConfig) (string, error) {
	var volume *disk.PartitionStat

	mounts, err := disk.PartitionsWithContext(context.Background(), true)
	if err != nil {
		return "", err
	}

MountsLoop:
	for _, m := range mounts {
		switch {
		case cfg.MountPath != nil:
			if *cfg.MountPath == m.Mountpoint {
				volume = &m
				break MountsLoop
			}
		case cfg.Device != nil:
			if *cfg.Device == m.Device {
				volume = &m
				break MountsLoop
			}
		case cfg.Name != nil:
			// warning: hacky
			// TODO: make less hacky
			volumeName := parseVolumeName(m.Mountpoint)
			if *cfg.Name == volumeName {
				volume = &m
				break MountsLoop
			}
		}
	}

	if volume != nil {
		return (*volume).Mountpoint, nil
	}

	return "", errors.New("no volumes matching the provided volume config found")
}

func parseVolumeName(mountPoint string) string {
	splitPath := strings.Split(mountPoint, "/")

	return splitPath[len(splitPath)-1]
}
