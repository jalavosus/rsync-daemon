package rsync

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/jalavosus/rsync-daemon/pkg/config"
	"github.com/jalavosus/rsync-daemon/pkg/volumes"
	log "github.com/sirupsen/logrus"
	"go-darwin.dev/hdiutil"
)

const (
	testVolumeSize    = 100 * 1024 * 1024 // 100MB
	defaultPermission = 0755
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}
}

func TestMain(m *testing.M) {
	var (
		volumeCleanupFunc func()
	)

	switch runtime.GOOS {
	case "darwin":
		volumeCleanupFunc = setupTestVolumeDarwin()
	default:
		volumeCleanupFunc = setupTestVolumeLinux()
	}

	defer volumeCleanupFunc()

	m.Run()
}

func TestSync(t *testing.T) {
	var backupVolume string
	switch runtime.GOOS {
	case "darwin":
		backupVolume = "/Volumes/rsyncd-test"
	default:
		backupVolume = "/tmp/rsyncd-test/mnt"
	}

	testCfg := &config.DaemonConfig{
		Include: []config.PathConfig{
			{
				Path:      "testdata/rsync",
				Recursive: true,
			},
		},
		BackupVolume: &volumes.VolumeConfig{
			MountPath: &backupVolume,
		},
	}

	testCfg.BackupDir, _ = volumes.GetVolumePath(*testCfg.BackupVolume)

	syncErr := Sync(testCfg)
	if syncErr != nil {
		t.Error(fmt.Errorf(string(syncErr.(*exec.ExitError).Stderr)))
	}
}

func setupTestVolumeDarwin() (cleanup func()) {
	diskImagePath := fmt.Sprintf("/tmp/rsync-daemon_rsync_test.img")
	tmpDisk, err := diskfs.Create(diskImagePath, testVolumeSize, diskfs.Raw)
	if err != nil {
		log.Fatal(err)
	}

	partitionSpec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeFat32,
		VolumeLabel: "rsyncd-test",
	}

	_, err = tmpDisk.CreateFilesystem(partitionSpec)
	if err != nil {
		log.Fatal(err)
	}

	deviceNode, mountErr := hdiutil.Attach(diskImagePath)
	if mountErr != nil {
		log.Fatal(mountErr)
	}

	cleanup = func() {
		unmountErr := hdiutil.Detach(deviceNode)
		if unmountErr != nil {
			log.Error(unmountErr)
			return
		}

		if err := os.Remove(diskImagePath); err != nil {
			log.Error(err)
		}

		return
	}

	return
}

func setupTestVolumeLinux() (cleanup func()) {
	ctx, ctxCancel := context.WithCancel(context.Background())

	tmpDir := "/tmp/rsyncd-test"

	tmp, err := ioutil.TempDir("", tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	backingDir := tmp + "/backing"
	mountPath := tmp + "/mnt"

	if err := os.Mkdir(backingDir, defaultPermission); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(mountPath, defaultPermission); err != nil {
		log.Fatal(err)
	}

	root := nodefs.NewMemNodeFSRoot(backingDir)
	conn := nodefs.NewFileSystemConnector(root, nil)

	server, err := fuse.NewServer(conn.RawFS(), mountPath, &fuse.MountOptions{
		Debug: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	go server.Serve()

	go func(ctx context.Context, server *fuse.Server) {
		for {
			select {
			case <-ctx.Done():
				if err := server.Unmount(); err != nil {
					log.Fatal(err)
				}
			default:
			}
		}
	}(ctx, server)

	cleanup = func() {
		ctxCancel()

		server.Wait()

		if err := os.RemoveAll(tmpDir); err != nil {
			log.Error(err)
		}
	}

	return
}
