package fsutils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
)

const (
	readable  string = "r"
	writeable string = "w"
)

func loadFileInfo(pathOrFileInfo interface{}) (os.FileInfo, error) {
	switch x := pathOrFileInfo.(type) {
	case string:
		return os.Stat(x)
	case os.FileInfo:
		return x, nil
	default:
		return nil, fmt.Errorf("param pathOrFileInfo must be one of (string, os.FileInfo); received: %[1]T", x)
	}
}

func parsePermissions(fileInfo os.FileInfo) (PathPermissions, error) {
	var pathPerms = PathPermissions{}

	filePerms := strings.Split(
		fileInfo.Mode().Perm().String(),
		"",
	)[1:]

	userUid, userGid, err := getUserUidGid()
	if err != nil {
		return pathPerms, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return pathPerms, errors.New("could not get deep stat information for file")
	}

	userOwnsFile := fmt.Sprintf("%[1]d", stat.Uid) == userUid
	fileInGroup := fmt.Sprintf("%[1]d", stat.Gid) == userGid

	switch {
	case userOwnsFile:
		pathPerms.Readable = filePerms[0] == readable
		pathPerms.Writeable = filePerms[1] == writeable
	case fileInGroup:
		pathPerms.Readable = filePerms[3] == readable
		pathPerms.Writeable = filePerms[4] == writeable
	}

	return pathPerms, nil
}

func getUserUidGid() (string, string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", err
	}

	return currentUser.Uid, currentUser.Gid, nil
}
