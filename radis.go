package main

import (
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
	i := sort.SearchStrings(len(g.Folders), artist)
	if i < len(g.Folders) && g.Folders[i] == artist {
		fmt.Println("Found artist ", artist, "in genre ", g.Name)
		return true
	}
	return false
}

//----------------

var reAlbum = regexp.MustCompile(`^([\p{L}\d_ ]+) \(([0-9]+)\) ([\p{L}\d_ ]+)$`)

type AlbumFolder struct {
	Root   string
	Path   string
	Artist string
	Year   int
	Title  string
	IsMP3  bool
}

func (a *AlbumFolder) IsAlbum() bool {
	// TODO isDIR + successful ExtractInfo
	return true
}

func (a *AlbumFolder) ExtractInfo() (err error) {
	matches := reAlbum.FindStringSubmatch(filepath.Base(a.Path))
	if matches {
		a.Artist = matches[1]
		a.Year = matches[2]
		a.Title = matches[3]
	}

	// TODO: IsMP3!!!!
	return
}

func (a *AlbumFolder) MoveToNewPath(root string, genre string) (err error) {
	if !a.IsAlbum() {
		fmt.Println("ERRRRRR")
		// TODO return ERR
	}

	directoryName, err := filepath.Rel(a.Root, a.Path)
	if err != nil {
		panic(err)
	}
	if a.IsMP3 {
		directoryName += " [MP3]"
	}
	newPath := filepath.Join(root, genre, a.Artist, directoryName)
	// comparer avec l'ancien
	if newPath != a.Path {
		fmt.Println(a.Path, " -> ", newPath)
		// if different, move folder
		os.Rename(a.Path, newPath)
	}
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

	re := regexp.MustCompile(`^([\p{L}\d_ ]+) \(([0-9]+)\) ([\p{L}\d_ ]+)$`)
	var artist, title string
	var year int

	err = filepath.Walk(root, func(path string, fileInfo os.FileInfo, _ error) error {
		relative, _ := filepath.Rel(root, path)
		fmt.Println("Scanning ", relative)

		if fileInfo.IsDir() {
			fmt.Println("-> is dir!")

			// TODO use AlbumFolder

			// find dir name
			directoryName := filepath.Base(relative)
			// extract artist, year, title, is_mp3
			matches := re.FindStringSubmatch(directoryName)
			if matches {
				artist = matches[1]
				year = matches[2]
				title = matches[3]

				found := false
				for genre, _ := range config {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(artist) {
						// prepare new path genre.Name/artist/directoryName
						newPath = filepath.Join(root, genre.Name, artist, directoryName)
						fmt.Println(newPath)

						// compare with current
						if newPath != path {
							fmt.Println(path, " -> ", newPath)
							// if different, move folder
							os.rename(path, newPath)
						}
						found = true
					}
				}
				// if not found, show error.
				if !found {
					fmt.Println("!!!! ERR artist not found: ", artist)

					// move to UNCATEGORIZED
					// prepare new path genre.Name/artist/directoryName
					newPath = filepath.Join(root, "UNCATEGORIZED", artist, directoryName)
					fmt.Println(newPath)

					// compare with current
					if newPath != path {
						fmt.Println(path, " -> ", newPath)
						// if different, move folder
						os.rename(path, newPath)
					}
				}
			} else {
				fmt.Println("Skipping ", directoryName)
			}

		}
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
