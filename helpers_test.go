package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	current, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get current directory!")
	}
	isEmpty, err := IsEmpty(current)
	if isEmpty || err != nil {
		t.Errorf("Current directory is not empty!")
	}
}

var errDirectoryDoesNotExist = errors.New("Path /ddsdcisj does not exist!!!")
var goPath = os.Getenv("GOPATH")

var testPaths = []struct {
	path         string
	existingPath string
	err          error
}{
	{"/tmp", "/tmp", nil},
	{"/ddsdcisj", "/ddsdcisj", errDirectoryDoesNotExist},
	{"/ddsdcisj", "/ddsdcisj", errDirectoryDoesNotExist},
	{"relative", filepath.Join(goPath, "src/github.com/barsanuphe/radis/relative"), errors.New("Path relative does not exist!!!")},
}

func TestGetExistingPath(t *testing.T) {
	for _, tp := range testPaths {
		path, err := GetExistingPath(tp.path)
		if path != tp.existingPath {
			t.Errorf("GetExistingPath(%s) returned %s, expected %s", tp.path, path, tp.existingPath)
		} else if err != nil && tp.err != nil && err.Error() != tp.err.Error() {
			t.Errorf("GetExistingPath(%s) returned err %s, expected %s", tp.path, err.Error(), tp.err.Error())
		}
	}
}

func TestHasNonFlacFiles(t *testing.T) {
	// TODO create fake directory with flac files
	current, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get current directory!")
	}
	hasNonFlac, err := HasNonFlacFiles(current)
	if !hasNonFlac || err != nil {
		t.Errorf("Current directory contains forbidden files!")
	}
}
