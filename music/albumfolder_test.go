package music

import (
	"errors"
	"testing"

	"fmt"
	"os"

	"github.com/barsanuphe/radis/config"
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
			t.Errorf("String(%s) returned %s, expected %s!", ta.folder, v, ta.expected)
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
	Folder          string
	Result          AlbumFolder
	Err             error
	ExpectedNewPath string
	ErrNewPath      error
	HasGenre        bool
}{
	{
		"music",
		AlbumFolder{Root: ".", Path: "music"},
		errors.New("Not an album!"),
		"",
		errors.New("Not an album!"),
		false,
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
		"genre1/PPP/arthi東京?-4. (2000) jqojdoijd(??)--+",
		nil,
		true,
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
		"UNCATEGORIZED/arthi/arthi (2000) jqojdoijd",
		nil,
		false,
	},
}

var c = config.Config{
	Paths: config.Paths{UnsortedSubdir: "UNCATEGORIZED"},
	Aliases: config.Aliases{
		config.Artist{MainAlias: "PPP", Aliases: []string{"arthi東京?-4."}},
		config.Artist{MainAlias: "CCC", Aliases: []string{"arthij"}},
	},
	Genres: config.Genres{
		config.Genre{Name: "genre1", Artists: []string{"PPP", "RRR"}},
	},
}

func TestExtractInfo(t *testing.T) {
	for _, ta := range albumsInfos {
		a := AlbumFolder{Root: ".", Path: ta.Folder}
		err := a.ExtractInfo()
		if err != ta.Err && a != ta.Result {
			t.Errorf("ExtractInfo(%s) returned %s, expected %s", ta.Folder, a.String(), ta.Result.String())
		}
	}
}

func TestFindNewPath(t *testing.T) {
	for _, ta := range albumsInfos {
		hasGenre, err := ta.Result.FindNewPath(c)
		if err != nil && err.Error() != ta.ErrNewPath.Error() {
			t.Errorf("TestFindNewPath(%s) returned err %s, expected %s", ta.Folder, err.Error(), ta.ErrNewPath.Error())
		}
		if hasGenre != ta.HasGenre {
			t.Errorf("TestFindNewPath(%s) returned hasGenre %v, expected %v", ta.Folder, hasGenre, ta.HasGenre)
		}
		if ta.Result.NewPath != ta.ExpectedNewPath {
			t.Errorf("TestFindNewPath(%s) returned NewPath %s, expected %s", ta.Folder, ta.Result.NewPath, ta.ExpectedNewPath)
		}
	}
}

func TestHasNonFlacFiles(t *testing.T) {
	// TODO create fake directory with flac files
	for _, ta := range albumsInfos {
		_, err := ta.Result.FindNewPath(c)
		// only testing on correct album folders
		if err == nil {
			hasNonFlac, err := ta.Result.HasNonFlacFiles()
			if os.IsNotExist(err) {
				fmt.Println(ta.Result.NewPath + " does not exist...")
			} else if !hasNonFlac || err != nil {
				t.Errorf("Directory " + ta.Result.NewPath + " contains forbidden files!")
			}
		}
	}
}

// TODO how to test MoveToNewPath????
