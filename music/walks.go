package music

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/barsanuphe/radis/config"
	"github.com/barsanuphe/radis/directory"
	"github.com/ttacon/chalk"
)

// TimeTrack can be used to evaluate the time spent in a function.
// usage: defer timeTrack(startTime) at the beginning of the function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

// SortAlbums scans the music collection root and reorders albums according to the configuration files.
func SortAlbums(c config.Config, doNothing bool) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	movedAlbums := 0
	uncategorized := 0
	foundAlbums := 0
	newAlbums := 0
	mp3Albums := 0

	dailyPlaylist, monthlyPlaylist := loadCurrentPlaylists(c)

	fmt.Printf("%sScanning for albums in %s...\n\n%s", chalk.Blue, c.Paths.Root, chalk.Reset)
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			a := Album{Root: c.Paths.Root, Path: path}
			if a.IsValidAlbum() {
				foundAlbums++
				if a.IsMP3 {
					mp3Albums++
				}
				hasGenre, err := a.FindNewPath(c)
				if err != nil {
					panic(err)
				}
				if !hasGenre {
					uncategorized++
				}

				originalRelative, _ := filepath.Rel(a.Root, a.Path)
				destRelative, _ := filepath.Rel(a.Root, a.NewPath)

				hasMoved, err := a.MoveToNewPath(doNothing)
				if err != nil {
					fmt.Println(chalk.Bold.TextStyle(chalk.Red.Color("!!! ERROR MOVING " + a.String())))
					fmt.Println(chalk.Bold.TextStyle(chalk.Red.Color("!!!\t    " + originalRelative + "\n!!!\t -> " + destRelative)))
				}
				if hasMoved {
					fmt.Println(chalk.Yellow.Color("+ " + a.String()))
					fmt.Println("\t    " + originalRelative + "\n\t -> " + destRelative)
					movedAlbums++
				}
				if a.IsNew(c) {
					// add to playlist automatically,
					fmt.Printf("%s\t    Adding to playlist.\n%s", chalk.Green, chalk.Reset)
					newAlbums++
					dailyPlaylist.contents = append(dailyPlaylist.contents, a)
					monthlyPlaylist.contents = append(monthlyPlaylist.contents, a)
				}
			}
		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Println(chalk.Blue)
	fmt.Printf("\n### Found %d albums including %d MP3 albums and %d new albums\n", foundAlbums, mp3Albums, newAlbums)
	if doNothing {
		fmt.Printf("### Sync would move %d albums.\n", movedAlbums)
	} else {
		fmt.Printf("### Moved %d albums.\n", movedAlbums)
	}
	if uncategorized != 0 {
		fmt.Println(chalk.Bold.TextStyle(chalk.Red.Color("\n!!!\n!!! " + strconv.Itoa(uncategorized) + " albums are still UNCATEGORIZED !!!\n!!!\n\n")))
	}
	if !doNothing {
		if err := writeCurrentPlaylists(dailyPlaylist, monthlyPlaylist); err != nil {
			panic(err)
		}
	}
	return
}

// FindNonFlacAlbums scan the music collection root and lists all albums of mp3 files instead of flac.
func FindNonFlacAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Printf("Scanning for non-Flac albums in %s.\n", c.Paths.Root)
	unFlagged := 0
	nonFlacAlbums := 0
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := Album{Root: c.Paths.Root, Path: path}
			if af.IsValidAlbum() {
				// scan contents for non-flac
				isNonFlac, err := af.HasNonFlacFiles()
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
	fmt.Printf("\n### Found %d non-Flac albums, including %d incorrectly flagged.\n", nonFlacAlbums, unFlagged)
	if unFlagged != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", unFlagged)
	}
	return
}

// DeleteEmptyFolders deletes empty folders that may appear after sorting albums.
func DeleteEmptyFolders(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Printf("Scanning for empty directories.\n\n")
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
				isEmpty, err := directory.IsEmpty(path)
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

	fmt.Printf("\n### Removed %d albums.\n", deletedDirectories)
	return
}
