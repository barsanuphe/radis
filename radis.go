// Radis is a tool to keep your music collection in great shape.
// 	see github.com/barsanuphe/radis
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/barsanuphe/radis/config"
	"github.com/codegangsta/cli"
)

// SORT ------------------------------------------------------------------------

// sortAlbums scans the music collection root and reorders albums according to the configuration files.
func sortAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for albums in " + c.Paths.Root + ".")
	movedAlbums := 0
	uncategorized := 0
	foundAlbums := 0
	mp3Albums := 0
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: c.Paths.Root, Path: path}
			if af.IsAlbum() {
				foundAlbums++
				if af.IsMP3 {
					mp3Albums++
				}
				hasMoved := false
				// fmt.Println("+ Found album: ", af.String())
				found := false

				// see if artist has known alias
				for _, alias := range c.Aliases {
					if alias.HasAlias(af.Artist) {
						af.MainAlias = alias.MainAlias
						break
					}
				}
				// find which genre the artist or main alias belongs to
				for _, genre := range c.Genres {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(af.MainAlias) || genre.HasCompilation(af.Title) {
						hasMoved, err = af.MoveToNewPath(genre.Name)
						found = true
						break
					}
				}
				if !found {
					uncategorized++
					hasMoved, err = af.MoveToNewPath(c.Paths.UnsortedSubdir)
				}
				if hasMoved {
					movedAlbums++
				}

				// TODO: detect if inside c.MainConfig.IncomingSubdir
				// try to find filepath.Rel( c.Paths.Root + c.Paths.IncomingSubdir, path), if err != nil, it's inside
				// TODO: if it is, add to playlist automatically
				// appendPlaylists(...)
				// NOTE: how to detect if the albums move again afterwards? need to update playlists?
			}
		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Printf("Found %d albums (%d MP3 albums), Moved %d.\n", foundAlbums, mp3Albums, movedAlbums)
	if uncategorized != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", uncategorized)
	}
	return
}

// LIST NON FLAC ---------------------------------------------------------------

// findNonFlacAlbums scan the music collection root and lists all albums of mp3 files instead of flac.
func findNonFlacAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for non-Flac albums in " + c.Paths.Root + ".")
	unFlagged := 0
	nonFlacAlbums := 0
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: c.Paths.Root, Path: path}
			if af.IsAlbum() {
				// scan contents for non-flac
				isNonFlac, err := HasNonFlacFiles(path)
				if err != nil {
					panic(err)
				}
				relativePath, _ := filepath.Rel(c.Paths.Root, path)
				if isNonFlac {
					fmt.Println("- ", relativePath)
					nonFlacAlbums++
				}
				if isNonFlac && !af.IsMP3 {
					unFlagged++
					fmt.Println("!!! ", relativePath, " not flagged as non FLAC!!!")
				}
				// NOTE: find falsely tagged folders? is that a thing?
			}
		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Printf("Found %d non-Flac albums, including %d incorrectly flagged.\n", nonFlacAlbums, unFlagged)
	if unFlagged != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", unFlagged)
	}
	return
}

// CLEAN -----------------------------------------------------------------------

// deleteEmptyFolders deletes empty folders that may appear after sorting albums.
func deleteEmptyFolders(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for empty directories.")
	deletedDirectories := 0
	deletedDirectoriesThisTime := 0
	atLeastOnce := false

	// loops until all levels of empty directories are deleted
	for !atLeastOnce || deletedDirectoriesThisTime != 0 {
		atLeastOnce = true
		deletedDirectoriesThisTime = 0
		err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
			// when an album has just been removed, Walk goes through it a second
			// time with an "file does not exist" error
			if os.IsNotExist(walkError) {
				return
			}
			if fileInfo.IsDir() {
				isEmpty, err := IsEmpty(path)
				if err != nil {
					panic(err)
				}
				if isEmpty {
					fmt.Println("Removing empty directory ", path)
					if err := os.Remove(path); err == nil {
						deletedDirectories++
						deletedDirectoriesThisTime++
					}
				}
			}
			return
		})
		if err != nil {
			fmt.Printf("Error!")
		}
	}

	fmt.Printf("Removed %d albums.\n", deletedDirectories)
	return
}

// MAIN ------------------------------------------------------------------------

func main() {
	fmt.Println("\n\tR A D I S\n\t---------\n")

	// load config
	rc := config.Config{}
	if err := rc.Load(); err != nil {
		panic(err)
	}
	// check config
	if err := rc.Check(); err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "R A D I S"
	app.Usage = "Organize your music collection."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"ls"},
			Usage:   "show configuration",
			Action: func(c *cli.Context) {
				// print config
				fmt.Println(rc.String())
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "sync folder according to configuration",
			Action: func(c *cli.Context) {
				// sort albums
				if err := sortAlbums(rc); err != nil {
					panic(err)
				}
				// scan again to remove empty directories
				if err := deleteEmptyFolders(rc); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:    "check",
			Aliases: []string{"find_awfulness"},
			Usage:   "check every album is a flac version, list the heretics.",
			Action: func(c *cli.Context) {
				// list non Flac albums
				if err := findNonFlacAlbums(rc); err != nil {
					panic(err)
				}
			},
		},
	}

	app.Run(os.Args)

	// write ordered config files
	if err := rc.Write(); err != nil {
		panic(err)
	}
}
