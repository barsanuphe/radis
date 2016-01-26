package music

import (
	"errors"
	"testing"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/barsanuphe/radis/config"
)

var albumsInfos = []struct {
	Folder          string
	Result          Album
	Err             error
	ExpectedNewPath string
	ErrNewPath      error
	HasGenre        bool
	HasNonFlac      bool
}{
	{
		"music",
		Album{Root: "/tmp/radis_test", Path: "/tmp/radis_test/music"},
		errors.New("Not an album!"),
		"",
		errors.New("Not an album!"),
		false,
		false,
	},
	{
		"arthi東京?-4. (2000) jqojdoijd(??)--+",
		Album{
			Root:      "/tmp/radis_test",
			Path:      "/tmp/radis_test/arthi東京?-4. (2000) jqojdoijd(??)--+",
			artist:    "arthi東京?-4.",
			mainAlias: "arthi東京?-4.",
			year:      "2000",
			title:     "jqojdoijd(??)--+",
			IsMP3:     false,
		},
		nil,
		"/tmp/radis_test/genre1/PPP/arthi東京?-4. (2000) jqojdoijd(??)--+",
		nil,
		true,
		false,
	},
	{
		"arthi (2000) jqojdoijd [MP3]",
		Album{
			Root:      "/tmp/radis_test",
			Path:      "/tmp/radis_test/arthi (2000) jqojdoijd",
			artist:    "arthi",
			mainAlias: "arthi",
			year:      "2000",
			title:     "jqojdoijd",
			IsMP3:     true,
		},
		nil,
		"/tmp/radis_test/UNCATEGORIZED/arthi/arthi (2000) jqojdoijd",
		nil,
		false,
		false,
	},
	{
		"artist (2000) title2",
		Album{
			Root:      "/tmp/radis_test",
			Path:      "/tmp/radis_test/UNCATEGORIZED/artist (2000) title2",
			artist:    "artist",
			mainAlias: "artist",
			year:      "2000",
			title:     "title2",
			IsMP3:     false,
		},
		nil,
		"/tmp/radis_test/UNCATEGORIZED/artist/artist (2000) title2",
		nil,
		false,
		true,
	},
}

var c = config.Config{
	Paths: config.Paths{Root: "/tmp/radis_test", UnsortedSubdir: "UNCATEGORIZED", IncomingSubdir: "INCOMING"},
	Aliases: config.Aliases{
		config.Artist{MainAlias: "PPP", Aliases: []string{"arthi東京?-4."}},
		config.Artist{MainAlias: "CCC", Aliases: []string{"arthij"}},
	},
	Genres: config.Genres{
		config.Genre{Name: "genre1", Artists: []string{"PPP", "RRR"}},
	},
}

func createTestFiles(c config.Config) {
	fmt.Println("Creating test files in ", c.Paths.Root)
	// create c.Paths.Root
	if err := os.MkdirAll(c.Paths.Root, 0777); err != nil {
		panic(err)
	}

	// create folders
	a0 := filepath.Join(c.Paths.Root, "arthi (2000) jqojdoijd")
	a1 := filepath.Join(c.Paths.Root, "arthi東京?-4. (2000) jqojdoijd(??)--+")
	a2 := filepath.Join(c.Paths.Root, "INCOMING", "artist (2000) title2")
	a3 := filepath.Join(c.Paths.Root, "UNCATEGORIZED", "artist (2000) title2")
	a4 := filepath.Join(c.Paths.Root, "genre1", "artist", "artist (2001) title3")
	for _, directory := range []string{a0, a1, a2, a3, a4} {
		if err := os.MkdirAll(directory, 0777); err != nil {
			panic(err)
		}
	}

	// create a few files
	f1 := filepath.Join(a1, "test.flac")
	f2 := filepath.Join(a2, "test.mp3")
	f3 := filepath.Join(a3, "test.wma")
	f4 := filepath.Join(a3, "test.flac")
	f5 := filepath.Join(a4, "test.flac")
	for _, file := range []string{f1, f2, f3, f4, f5} {
		if err := ioutil.WriteFile(file, []byte{}, 0777); err != nil {
			panic(err)
		}
	}
}

func cleanTestFiles(c config.Config) {
	fmt.Println("Cleaning ", c.Paths.Root)
	if err := os.RemoveAll(c.Paths.Root); err != nil {
		panic(err)
	}
}

// TestMain runs all tests, after creating temporary test files.
func TestMain(m *testing.M) {
	createTestFiles(c)
	result := m.Run()
	cleanTestFiles(c)
	os.Exit(result)
}

var albumsFolders = []struct {
	folder         string
	expectedString string
	isAlbum        bool
}{
	{"hop", "/ () ", false},
	{"arthi (2000) jqojdoijd", "arthi/arthi (2000) jqojdoijd", true},
	{"arthi (2000) jqojdoijd [MP3]", "arthi/arthi (2000) jqojdoijd [MP3]", true},
	{"arthi (2000) jqojdoijd [EP]", "arthi/arthi (2000) jqojdoijd [EP]", true},
	{"arthi (20010) jqojdoijd [EP]", "/ () ", false},
	{"arthi (2010) jqojdoijd (??ï4é)--+", "arthi/arthi (2010) jqojdoijd (??ï4é)--+", true},
}

func TestString(t *testing.T) {
	for _, ta := range albumsFolders {
		a := Album{Root: ".", Path: ta.folder}
		a.extractInfo()
		if v := a.String(); v != ta.expectedString {
			t.Errorf("String(%s) returned %s, expected %s!", ta.folder, v, ta.expectedString)
		}
	}
}

func TestIsValidAlbum(t *testing.T) {
	for _, ta := range albumsFolders {
		a := Album{Root: ".", Path: ta.folder}
		v := a.IsValidAlbum()
		if v != ta.isAlbum {
			t.Errorf("IsAlbum(%s) returned %v, expected %v", ta.folder, v, ta.isAlbum)
		}
		// should return true the second time
		if v && !a.IsValidAlbum() {
			t.Errorf("IsAlbum(%s) returned %v, expected %v", ta.folder, v, ta.isAlbum)
		}
	}
}

func TestExtractInfo(t *testing.T) {
	for _, ta := range albumsInfos {
		a := Album{Root: c.Paths.Root, Path: ta.Folder}
		err := a.extractInfo()
		if err != ta.Err && a.String() != ta.Result.String() {
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
	for _, ta := range albumsInfos {
		_, err := ta.Result.FindNewPath(c)
		// only testing on correct album folders
		if err == nil {
			hasNonFlac, err := ta.Result.HasNonFlacFiles()
			if os.IsNotExist(err) {
				t.Errorf(ta.Result.NewPath + " does not exist...")
			}
			if hasNonFlac != ta.HasNonFlac {
				t.Errorf("HasNonFlacFiles(%s) returned %v, expected %v", ta.Folder, hasNonFlac, ta.HasNonFlac)
			}
			if err != nil {
				t.Errorf("HasNonFlacFiles(%s) returned en error!: %s", err.Error())
			}
		}
	}
}

// TODO MoveToNewPath, GetMusicFiles
