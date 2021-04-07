package fsutils

import (
	"os"
	"path"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	testFilesDir    string = "testdata/fsutils"
	userOwnedPrefix string = "user_owned_"
	rootOwnedPrefix string = "root_owned_"
)

const (
	readOnly            string = "ro"
	readWrite           string = "rw"
	rootOnlySuffix      string = "root_only.txt"
	userReadOnlySuffix         = "user_" + readOnly + ".txt"
	userReadWriteSuffix        = "user_" + readWrite + ".txt"
)

var (
	rootOnlyTestFile = path.Join(testFilesDir, rootOwnedPrefix+rootOnlySuffix)
	rootOnlyTestDir  = path.Join(testFilesDir, rootOwnedPrefix+"dir")
)

type pathPermissionsTest struct {
	path              string
	shouldBeReadable  bool
	shouldBeWriteable bool
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}
}

func TestGetPathPermissions(t *testing.T) {
	testCases := []pathPermissionsTest{
		{
			path:              path.Join(testFilesDir, rootOwnedPrefix+userReadOnlySuffix),
			shouldBeReadable:  true,
			shouldBeWriteable: false,
		},
		{
			path:              path.Join(testFilesDir, rootOwnedPrefix+userReadWriteSuffix),
			shouldBeReadable:  true,
			shouldBeWriteable: true,
		},
		{
			path:              path.Join(testFilesDir, userOwnedPrefix+userReadOnlySuffix),
			shouldBeReadable:  true,
			shouldBeWriteable: false,
		},
		{
			path:              path.Join(testFilesDir, userOwnedPrefix+userReadWriteSuffix),
			shouldBeReadable:  true,
			shouldBeWriteable: true,
		},
	}

	if os.Getuid() == 0 {
		os.Chown(rootOnlyTestFile, 0, 0)
		os.Chmod(rootOnlyTestFile, 0600)
		os.Chown(rootOnlyTestDir, 0, 0)
		os.Chmod(rootOnlyTestDir, 0600)

		testCases = append(testCases,
			pathPermissionsTest{
				path:              path.Join(testFilesDir, rootOwnedPrefix+rootOnlySuffix),
				shouldBeReadable:  false,
				shouldBeWriteable: false,
			},
			pathPermissionsTest{
				path:              path.Join(testFilesDir, rootOwnedPrefix+"dir"),
				shouldBeReadable:  false,
				shouldBeWriteable: false,
			})
	}

	for _, tc := range testCases {
		pathPerms, err := GetPathPermissions(tc.path)
		if err != nil {
			t.Log(err)
			t.Fail()
			continue
		}

		assert.Equal(t, tc.shouldBeReadable, pathPerms.Readable,
			"expected path %[1]s to have Readable value of %[2]t; got %[3]t",
			tc.path,
			tc.shouldBeReadable,
			pathPerms.Readable,
		)

		assert.Equal(t, tc.shouldBeWriteable, pathPerms.Writeable,
			"expected path %[1]s to have Writeable value of %[2]t; got %[3]t",
			tc.path,
			tc.shouldBeWriteable,
			pathPerms.Writeable,
		)
	}
}
