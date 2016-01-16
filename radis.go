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

func sortFolders(root string, genres []Genre, aliases MainAlias) (err error) {
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

				// see if artist has known alias
				for _, alias := range aliases {
					if alias.HasAlias(af.Artist) {
						af.MainAlias = alias.MainAlias
						break
					}
				}
				// find which genre the artist or main alias belongs to
				for _, genre := range genres {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(af.MainAlias) {
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
	fmt.Println("\n\tR A D I S\n\t---------\n")
	pwd, _ := os.Getwd()
	aliasesConfigFile := filepath.Join(pwd, "radis_aliases.yaml")
	genresConfigFile := filepath.Join(pwd, "radis.yaml")

	// load config files
	aliases := MainAlias{}
	if err := aliases.Load(aliasesConfigFile); err != nil {
		panic(err)
	}
	genres := Config{}
	if err := genres.Load(genresConfigFile); err != nil {
		panic(err)
	}

	// print config
	//fmt.Println(aliases.String())
	//fmt.Println(genres.String())

	// scan folder in root
	root := filepath.Join(pwd, "test/")
	if err := sortFolders(root, genres, aliases); err != nil {
		panic(err)
	}
	// TODO scan again to remove empty directories

	// write ordered config files
	if err := aliases.Write(aliasesConfigFile); err != nil {
		panic(err)
	}
	if err := genres.Write(genresConfigFile); err != nil {
		panic(err)
	}
}
