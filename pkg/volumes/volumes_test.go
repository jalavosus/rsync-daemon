package volumes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pointers "github.com/xplorfin/pointerutils"
)

type testVolumeConfig struct {
	VolumeConfig
	shouldFail bool
}

func TestReadVolumes(t *testing.T) {
	testCfgs := []testVolumeConfig{
		{
			VolumeConfig: VolumeConfig{
				Name: pointers.FromString("Kitties Causing Fires"),
			},
			shouldFail: false,
		},
		{
			VolumeConfig: VolumeConfig{
				Device: pointers.FromString("/dev/disk1s1"),
			},
			shouldFail: false,
		},
		{
			VolumeConfig: VolumeConfig{
				MountPath: pointers.FromString("/Volumes/Kitties Causing Fires"),
			},
			shouldFail: false,
		},
		{
			VolumeConfig: VolumeConfig{
				MountPath: pointers.FromString("/Volumes/Test"),
			},
			shouldFail: true,
		},
	}

	for _, testCfg := range testCfgs {
		_, err := GetVolumePath(testCfg.VolumeConfig)

		if err != nil {
			if !assert.True(t, testCfg.shouldFail) {
				assert.Fail(t, err.Error())
			}
		}
	}
}
