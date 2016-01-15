package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

// TODO : save config with sorted Folders

//----------------
type Genre struct {
	Name    string
	Folders []string
}

func (g *Genre) String() string {
	txt := g.Name + ":\n"
	for _, artist := range g.Folders {
		txt += "\t- " + artist + "\n"
	}
	return txt
}

func (g *Genre) HasArtist(artist string) bool {
	sort.Strings(g.Folders)
	i := sort.SearchStrings(g.Folders, artist)
	if i < len(g.Folders) && g.Folders[i] == artist {
		// fmt.Println("++ Found artist ", artist, "in genre ", g.Name)
		return true
	}
	return false
}

//----------------

// TODO check if we need other unicode classes
var reAlbum = regexp.MustCompile(`^([\p{L}\d_ ]+) \(([0-9]+)\) ([\p{L}\d_ ]+)(\s\[MP3\])?$`)

type AlbumFolder struct {
	Root   string
	Path   string
	Artist string
	Year   string
	Title  string
	IsMP3  bool
}

func (a *AlbumFolder) String() string {
	if a.IsMP3 {
		return a.Artist + " (" + a.Year + ") " + a.Title + " [MP3]"
	} else {
		return a.Artist + " (" + a.Year + ") " + a.Title
	}
}

func (a *AlbumFolder) IsAlbum() bool {
	// TODO isDIR
	if err := a.ExtractInfo(); err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

func (a *AlbumFolder) ExtractInfo() (err error) {
	matches := reAlbum.FindStringSubmatch(filepath.Base(a.Path))
	if len(matches) > 0 {
		a.Artist = matches[1]
		a.Year = matches[2]
		a.Title = matches[3]
		a.IsMP3 = matches[4] != ""
	} else {
		err = errors.New("Not an album!")
	}
	// TODO: IsMP3!!!!
	return
}

func (a *AlbumFolder) MoveToNewPath(genre string) (err error) {
	// TODO: return bool hasMoved

	if !a.IsAlbum() {
		return errors.New("Cannot move, not an album.")
	}

	directoryName := filepath.Base(a.Path)
	newPath := filepath.Join(a.Root, genre, a.Artist, directoryName)
	// comparer avec l'ancien
	if newPath != a.Path {
		// if different, move folder
		originalRelative, _ := filepath.Rel(a.Root, a.Path)
		destRelative, _ := filepath.Rel(a.Root, newPath)
		fmt.Println("+ "+originalRelative, " -> ", destRelative)

		// TODO create newPath parent if it does not exist
		newPathParent := filepath.Dir(newPath)
		if _, err = os.Stat(newPathParent); os.IsNotExist(err) {
			// newPathParent does not exist, creating
			err = os.MkdirAll(newPathParent, 0777)
			if err != nil {
				panic(err)
			}
		}

		// move
		err = os.Rename(a.Path, newPath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

//----------------

func printConfig(config []Genre) {
	for _, genre := range config {
		fmt.Println(genre.String())
	}
}

func readConfig(path string) (Config []Genre, err error) {
	data, err := ioutil.ReadFile("radis.yaml")
	if err != nil {
		panic(err)
	}

	m := make(map[string][]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}

	for genre := range m {
		var newGenre Genre
		newGenre.Name = genre
		newGenre.Folders = m[genre]
		Config = append(Config, newGenre)
	}
	return
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

func sortFolders(root string, config []Genre) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	err = filepath.Walk(root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}


		if fileInfo.IsDir() {
			// relative, _ := filepath.Rel(root, path)
			// fmt.Println("Scanning ", relative)
			af := AlbumFolder{Root: root, Path: path}
			if af.IsAlbum() {
				fmt.Println("+ Found album: ", af.String())
				found := false
				for _, genre := range config {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(af.Artist) {
						err = af.MoveToNewPath(genre.Name)
						found = true
						break
					}
				}
				if !found {
					err = af.MoveToNewPath("UNCATEGORIZED")
				}
			} else {
				// fmt.Println("++ Skipping, not an album.")
			}

		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Printf("\rScanning: Done.\n")
	return
}

//----------------

func main() {
	fmt.Println("R A D I S\n---------\n")
	pwd, _ := os.Getwd()

	config, err := readConfig(filepath.Join(pwd, "radis.yaml"))
	if err != nil {
		panic(err)
	}

	printConfig(config)

	// scan folder in root
	root := filepath.Join(pwd, "test/")
	err = sortFolders(root, config)
	if err != nil {
		panic(err)
	}
}
