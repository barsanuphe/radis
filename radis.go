package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

func readConfig(path string) (genres AllGenres, err error) {
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
		genres = append(genres, newGenre)
	}
	return
}

func writeConfig(config AllGenres, path string) (err error) {
	sort.Sort(config)
	m := make(map[string][]string)
	for _, genre := range config {
		sort.Strings(genre.Folders)
		m[genre.Name] = genre.Folders
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Println("error: %v", err)
		panic(err)
	}
	err = ioutil.WriteFile(path, d, 0777)
	if err != nil {
		fmt.Println("error: %v", err)
	}
	return
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

func sortFolders(root string, config []Genre) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for albums.")
	movedAlbums := 0
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
				hasMoved := false
				fmt.Println("+ Found album: ", af.String())
				found := false
				for _, genre := range config {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(af.Artist) {
						hasMoved, err = af.MoveToNewPath(genre.Name)
						found = true
						break
					}
				}
				if !found {
					hasMoved, err = af.MoveToNewPath("UNCATEGORIZED")
				}
				if hasMoved {
					movedAlbums++
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
	fmt.Printf("Moved %d albums.\n", movedAlbums)
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

	fmt.Println(config.String())

	// scan folder in root
	root := filepath.Join(pwd, "test/")
	err = sortFolders(root, config)
	if err != nil {
		panic(err)
	}

	err = writeConfig(config, filepath.Join(pwd, "radis_out.yaml"))
	if err != nil {
		panic(err)
	}
}

