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
	{"arthi (2000) jqojdoijd [EP]", "arthi/arthi (2000) jqojdoijd [EP]"},
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
	{"arthi (2000) jqojdoijd [EP]", true},
	{"arthi (20010) jqojdoijd", false},
	{"arthi (2010) jqojdoijd (??)--+", true},
}

func TestIsAlbum(t *testing.T) {
	for _, ta := range albumsPaths {
		a := AlbumFolder{Root: ".", Path: ta.folder}
		v := a.IsAlbum()
		if v != ta.expected {
			t.Errorf("IsAlbum(%s) returned %v, expected %v", ta.folder, v, ta.expected)
		}
		// should return true the second time
		if v && !a.IsAlbum() {
			t.Errorf("IsAlbum(%s) returned %v, expected %v", ta.folder, v, ta.expected)
		}
	}
}

var albumsInfos = []struct {
	Folder string
	Result AlbumFolder
	Err    error
}{
	{
		"hop",
		AlbumFolder{Root: ".", Path: "hop"},
		errors.New("Not an album!"),
	},
	{
		"arthi東京?-4. (2000) jqojdoijd(??)--+",
		AlbumFolder{
			Root:      ".",
			Path:      "arthi東京?-4. (2000) jqojdoijd(??)--+",
			Artist:    "arthi東京?-4.",
			MainAlias: "arthi東京?-4.",
			Year:      "2000",
			Title:     "jqojdoijd(??)--+",
			IsMP3:     false,
		},
		nil,
	},
	{
		"arthi (2000) jqojdoijd [MP3]",
		AlbumFolder{
			Root:      ".",
			Path:      "arthi (2000) jqojdoijd",
			Artist:    "arthi",
			MainAlias: "arthi",
			Year:      "2000",
			Title:     "jqojdoijd",
			IsMP3:     true,
		},
		nil,
	},
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
