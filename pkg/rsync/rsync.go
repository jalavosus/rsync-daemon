package rsync

import (
	"os/exec"

	"github.com/jalavosus/rsync-daemon/pkg/config"
)

const cmd string = "rsync"

var defaultRsyncArgs = []string{
	"--ignore-existing",
	"--update",
	"--prune-empty-dirs",
}

func Sync(cfg *config.DaemonConfig) error {
	for _, include := range cfg.Include {
		var rsyncArgs = defaultRsyncArgs

		if include.Recursive {
			rsyncArgs = append(rsyncArgs, "--recursive")
		}

		rsyncArgs = append(rsyncArgs, include.Path)

		rsyncArgs = append(rsyncArgs, cfg.BackupDir)

		command := exec.Command(cmd, rsyncArgs...)

		_, err := command.Output()
		if err != nil {
			return err
		}
	}

	return nil
}
