package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/jalavosus/rsync-daemon/pkg/fsutils"
	"github.com/jalavosus/rsync-daemon/pkg/volumes"
	log "github.com/sirupsen/logrus"
)

func ValidateConfig(config *DaemonConfig) error {
	var hasPathValidationErrors bool

	if config.BackupVolume == nil {
		return errors.New("must provide a configuration for the backup volume to use")
	}

	if _, err := volumes.GetVolumePath(*config.BackupVolume); err != nil {
		return errors.New("error parsing volume backup volume configuration: " + err.Error())
	}

	if len(config.Include) == 0 {
		return errors.New("field `include` cannot be empty")
	}

	includesValidationErrors := validatePaths(config.Include)

	if len(includesValidationErrors) > 0 {
		hasPathValidationErrors = true

		log.Warnf("----Validation errors in config.Include----")
		for _, valErr := range includesValidationErrors {
			log.Warn(valErr)
		}
	}

	if hasPathValidationErrors {
		return errors.New("validation errors were found in config.Include; see above output for details")
	}

	return nil
}

func validatePaths(paths []PathConfig) []error {
	var pathValidationErrors []error

	for _, p := range paths {
		pathInfo, err := os.Stat(p.Path)
		if err != nil {
			continue
		}

		pathExists, err := fsutils.CheckPathExists(pathInfo)
		if err != nil {
			pathValidationErrors = append(pathValidationErrors, fmt.Errorf("os error: %[1]v", err))
			continue
		}

		if !pathExists {
			pathValidationErrors = append(
				pathValidationErrors,
				fmt.Errorf("path %[1]s does not exist or isn't valid", p.Path),
			)
			continue
		}

		pathPermissions, err := fsutils.GetPathPermissions(pathInfo)
		if err != nil {
			pathValidationErrors = append(
				pathValidationErrors,
				fmt.Errorf("os error: %[1]v", err),
			)
			continue
		}

		if !pathPermissions.Readable {
			pathValidationErrors = append(
				pathValidationErrors,
				fmt.Errorf("path %[1]s is not readable by the current user or group", p.Path),
			)
		}

		if fsutils.PathIsDirectory(pathInfo) && !p.Recursive {
			pathValidationErrors = append(
				pathValidationErrors,
				fmt.Errorf("path %[1]s is a directory but `recursive` was set to false", p.Path),
			)
			continue
		}
	}

	return pathValidationErrors
}
