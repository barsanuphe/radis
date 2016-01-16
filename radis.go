package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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

	config := Config{}
	if err := config.Load(filepath.Join(pwd, "radis.yaml")); err != nil {
		panic(err)
	}
	fmt.Println(config.String())

	// scan folder in root
	root := filepath.Join(pwd, "test/")
	if err := sortFolders(root, config); err != nil {
		panic(err)
	}

	if err := config.Write(filepath.Join(pwd, "radis_out.yaml")); err != nil {
		panic(err)
	}
}
