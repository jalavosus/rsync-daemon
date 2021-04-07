package config

import (
	"os"
	"path"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}
}

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("testdata/test_cfg.yaml", true)
	assert.Nil(t, err)
}
