package music

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barsanuphe/radis/config"
	"github.com/barsanuphe/radis/directory"
)

// TimeTrack can be used to evaluate the time spent in a function.
// usage: defer timeTrack(startTime) at the beginning of the function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

// SortAlbums scans the music collection root and reorders albums according to the configuration files.
func SortAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	movedAlbums := 0
	uncategorized := 0
	foundAlbums := 0
	mp3Albums := 0

	dailyPlaylist, monthlyPlaylist := loadCurrentPlaylists(c)

	fmt.Println("Scanning for albums in " + c.Paths.Root + ".")
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
				hasGenre, err := af.FindNewPath(c)
				if err != nil {
					panic(err)
				}
				if !hasGenre {
					uncategorized++
				}
				hasMoved, err := af.MoveToNewPath()
				if err != nil {
					panic(err)
				}
				if hasMoved {
					movedAlbums++
				}

				// find out if we should append playlists
				isInsideIncoming := strings.Contains(path, filepath.Join(c.Paths.Root, c.Paths.IncomingSubdir))
				if isInsideIncoming {
					// album is inside INCOMING dir, add to playlist automatically
					fmt.Println("  ++ " + af.String() + " was in INCOMING, adding to playlist.")
					dailyPlaylist.contents = append(dailyPlaylist.contents, af)
					monthlyPlaylist.contents = append(monthlyPlaylist.contents, af)
				}
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

	if err := writeCurrentPlaylists(dailyPlaylist, monthlyPlaylist); err != nil {
		panic(err)
	}
	return
}

// FindNonFlacAlbums scan the music collection root and lists all albums of mp3 files instead of flac.
func FindNonFlacAlbums(c config.Config) (err error) {
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
	fmt.Printf("Found %d non-Flac albums, including %d incorrectly flagged.\n", nonFlacAlbums, unFlagged)
	if unFlagged != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", unFlagged)
	}
	return
}

// DeleteEmptyFolders deletes empty folders that may appear after sorting albums.
func DeleteEmptyFolders(c config.Config) (err error) {
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

	fmt.Printf("Removed %d albums.\n", deletedDirectories)
	return
}
