package fsutils

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

type PathPermissions struct {
	Readable  bool // Is this path readable by the current user?
	Writeable bool // Is this path writeable by the current user?
}

func GetPathPermissions(pathParam interface{}) (PathPermissions, error) {
	var pathPerms = PathPermissions{}

	pathInfo, err := loadFileInfo(pathParam)
	if err != nil {
		return pathPerms, err
	}

	pathPerms, err = parsePermissions(pathInfo)
	if err != nil {
		log.Fatal(err)
	}

	return pathPerms, nil
}

func CheckPathExists(pathParam interface{}) (bool, error) {
	if _, err := loadFileInfo(pathParam); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil // Why return an error if existence checking is the whole point of the function?
		}

		return false, err
	}

	return true, nil
}

func PathIsDirectory(pathParam interface{}) bool {
	pathInfo, _ := loadFileInfo(pathParam)

	return pathInfo.IsDir()
}
