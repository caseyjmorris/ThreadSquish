package testHelpers

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetFixturePath(fixture string) string {
	return path.Join(getRootDir(), "testFiles", fixture)
}

func getRootDir() string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "threadsquish") {
		wd = filepath.Dir(wd)
	}

	return wd
}
