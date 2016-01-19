package main

import (
	"errors"
	"testing"
)

var albumsStrings = []struct {
	folder   string
	expected string
}{
	{"hop", "/ () "},
	{"arthi (2000) jqojdoijd", "arthi/arthi (2000) jqojdoijd"},
	{"arthi (2000) jqojdoijd [MP3]", "arthi/arthi (2000) jqojdoijd [MP3]"},
	// TODO more
}

func TestString(t *testing.T) {
	for _, ta := range albumsStrings {
		a := AlbumFolder{Root: ".", Path: ta.folder}
		a.ExtractInfo()
		if v := a.String(); v != ta.expected {
			t.Errorf("String(%s) returned %s, expected %s", ta.folder, v, ta.expected)
		}
	}
}

var albumsPaths = []struct {
	folder   string
	expected bool
}{
	{"hop", false},
	{"arthi (2000) jqojdoijd", true},
	{"arthi (2000) jqojdoijd [MP3]", true},
	// TODO more
}

func TestIsAlbum(t *testing.T) {
	for _, ta := range albumsPaths {
		a := AlbumFolder{Root: ".", Path: ta.folder}
		if v := a.IsAlbum(); v != ta.expected {
			t.Errorf("IsAlbum(%s) returned %v, expected %v", ta.folder, v, ta.expected)
		}
	}
}

var albumsInfos = []struct {
	Folder string
	Result AlbumFolder
	Err    error
}{
	{"hop", AlbumFolder{Root:".", Path: "hop"}, errors.New("Not an album!")},
	{"arthi (2000) jqojdoijd", AlbumFolder{
		Root:      ".",
		Path:      "arthi (2000) jqojdoijd",
		Artist:    "arthi",
		MainAlias: "arthi",
		Year:      "2000",
		Title:     "jqojdoijd",
		IsMP3:     false,
	}, nil},
	{"arthi (2000) jqojdoijd [MP3]", AlbumFolder{
		Root:      ".",
		Path:      "arthi (2000) jqojdoijd",
		Artist:    "arthi",
		MainAlias: "arthi",
		Year:      "2000",
		Title:     "jqojdoijd",
		IsMP3:     true,
	}, nil},
	// TODO more
}

func TestExtractInfo(t *testing.T) {
	for _, ta := range albumsInfos {
		a := AlbumFolder{Root: ".", Path: ta.Folder}
		err := a.ExtractInfo()
		if err != ta.Err && a != ta.Result {
			// TODO print err too
			t.Errorf("ExtractInfo(%s) returned %d, expected %d", ta.Folder, a.String(), ta.Result)
		}
	}
}

// TODO how to test MoveToNewPath????
