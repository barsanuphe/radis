package main

import (
	"errors"
	"testing"
)

var testPaths = []struct {
	path         string
	existingPath string
	err          error
}{
	{"/tmp", "/tmp", nil},
	{"/ddsdcisj", "/ddsdcisj", errors.New("Directory /ddsdcisj does not exist!!!")},
	// TODO more avec cas en relatif
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
