package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jalavosus/rsync-daemon/pkg/fsutils"
	"github.com/jalavosus/rsync-daemon/pkg/volumes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"os"
	"path"
)

const DefaultConfigFilename string = ".rsyncdaemoncfg.yaml"

var (
	DefaultConfigDir  string
	DefaultConfigPath string
)

func init() {
	DefaultConfigDir, _ = os.UserHomeDir()
	DefaultConfigPath = path.Join(DefaultConfigDir, DefaultConfigFilename)
}

type DaemonConfig struct {
	// Files or directories to include in the backup.
	Include []PathConfig `yaml:"include,flow"`
	// Files or directories to specifically exclude from the backup.
	Exclude []PathConfig `yaml:"exclude,flow"`
	// External drive to back up files to.
	BackupVolume *volumes.VolumeConfig `yaml:"backup_volume,flow"`
	// Directory on the external volume to backup to; defaults to the root
	// directory of the volume.
	BackupDir string `yaml:"backup_dir"`
}

type PathConfig struct {
	Path      string `yaml:"path"`
	Recursive bool   `yaml:"recursive"`
}

func LoadConfig(configPath string, ignoreValidation bool) (*DaemonConfig, error) {
	loadedConfig, err := readConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}

	if !ignoreValidation {
		validationError := ValidateConfig(loadedConfig)
		if validationError != nil {
			log.Fatalf("error validating config: %[1]v", validationError)
		}
	}

	volPath, err := volumes.GetVolumePath(*loadedConfig.BackupVolume)
	if err != nil {
		if !ignoreValidation {
			return nil, err
		}
	}

	if loadedConfig.BackupDir == "" {
		loadedConfig.BackupDir = volPath
	} else {
		loadedConfig.BackupDir = filepath.Join(volPath, loadedConfig.BackupDir)
	}

	return loadedConfig, nil
}

func readConfigFromFile(configPath string) (*DaemonConfig, error) {
	var loadedConfig = &DaemonConfig{}

	// we already know that DefaultConfigPath is a fully qualified path
	if configPath != DefaultConfigPath {
		pathExists, err := fsutils.CheckPathExists(configPath)
		if err != nil {
			return nil, err // likely a fatal error somewhere
		}

		switch pathExists {
		case true:
			if fsutils.PathIsDirectory(configPath) {
				configFilePath := path.Join(configPath, DefaultConfigFilename)
				fileExistsInDirectory, err := fsutils.CheckPathExists(configFilePath)
				if err != nil {
					return nil, err
				}

				if !fileExistsInDirectory {
					log.Warnf("%[1]s does not exist in directory %[2]s; "+
						"defaulting to creating a new config file at %[3]s", DefaultConfigFilename, configPath, configFilePath)
				}

				configPath = configFilePath
			}
		default:
			log.Warnf("%[1]s is not a valid file or directory, "+
				"defaulting to creating a new config file at %[2]s", configPath, DefaultConfigPath)
			configPath = DefaultConfigPath
		}
	}

	// check that we have permission to write to our chosen directory
	filePerms, err := fsutils.GetPathPermissions(configPath)
	if err != nil {
		return nil, err
	}

	if !filePerms.Readable {
		return nil, fmt.Errorf("it looks like rsync-daemon doesn't have permission to read %[1]s; "+
			"please check the file permissions and try again", configPath)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var unmarshalFunc func([]byte, interface{}) error

	switch path.Ext(configPath) {
	case ".yaml":
		unmarshalFunc = yaml.Unmarshal
	case ".json":
		unmarshalFunc = json.Unmarshal
	default:
		return nil, fmt.Errorf("%[1]s is not a supported file extension for config", path.Ext(configPath))
	}

	if err := unmarshalFunc(configBytes, &loadedConfig); err != nil {
		return nil, err
	}

	return loadedConfig, nil
}
