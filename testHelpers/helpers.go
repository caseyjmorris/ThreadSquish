package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetFixturePath(fixture string) string {
	return filepath.Join(getRootDir(), "testFiles", fixture)
}

func Unix2Dos(path string) error {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error opening %q:  %s", path, err)
	}
	lineBreak := regexp.MustCompile("\r?\n")
	dosified := lineBreak.ReplaceAll(text, []byte("\r\n"))
	err = ioutil.WriteFile(path, dosified, 0644)

	if err != nil {
		return fmt.Errorf("error writing %q:  %s", path, err)
	}
	return nil
}

func getRootDir() string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "threadsquish") {
		wd = filepath.Dir(wd)
	}

	return wd
}
