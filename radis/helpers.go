package radis

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/barsanuphe/radis/config"
)

// timeTrack can be used to evaluate the time spent in a function.
// usage: defer timeTrack(startTime) at the beginning of the function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

// removeDuplicatePaths takes a slice of paths, return one without duplicates
func removeDuplicatePaths(a []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

// IsEmpty checks if a directory is empty.
func IsEmpty(directory string) (bool, error) {
	f, err := os.Open(directory)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// GetExistingPath ensures a path actually exists, and returns an existing absolute path or an error.
func GetExistingPath(path string) (existingPath string, err error) {
	// check root exists or pwd+root exists
	if filepath.IsAbs(path) {
		existingPath = path
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		existingPath = filepath.Join(pwd, path)
	}
	// check root exists
	if _, err = os.Stat(existingPath); os.IsNotExist(err) {
		err = errors.New("Path " + path + " does not exist!!!")
	}
	return
}

func getFiles(path string) (contents []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	// get contents
	fileList, err := f.Readdirnames(-1)
	if err == io.EOF {
		return []string{}, nil
	}
	return fileList, err
}

// HasNonFlacFiles returns true if an album contains files other than flac songs and cover pictures.
func HasNonFlacFiles(albumPath string) (bool, error) {
	fileList, err := getFiles(albumPath)
	if err != nil {
		return false, err
	}
	// check for suspicious files
	hasNonFlac := false
	for _, file := range fileList {
		switch filepath.Ext(file) {
		case ".flac", ".jpg", ".jpeg", ".png":
			// accepted extensions
		case ".mp3", ".wma", ".m4a":
			hasNonFlac = true
			break
		default:
			fmt.Println("Found suspicious file ", file, " in ", albumPath)
			hasNonFlac = true
			break
		}
	}
	return hasNonFlac, err
}

// GetMusicFiles returns a list of flacs and mp3s in an album
func GetMusicFiles(albumPath string) (contents []string, err error) {
	fileList, err := getFiles(albumPath)
	if err != nil {
		return []string{}, err
	}
	// check for music files
	for _, file := range fileList {
		switch filepath.Ext(file) {
		case ".flac", ".mp3":
			// accepted extensions
			contents = append(contents, filepath.Join(albumPath, file))
		}
	}
	sort.Strings(contents)
	return
}

// GetPlaylistFiles returns a list of .m3u files.
func GetPlaylistFiles(playlistRoot string) (contents []string, err error) {
	fileList, err := getFiles(playlistRoot)
	if err != nil {
		return []string{}, err
	}
	// check for m3u files
	for _, file := range fileList {
		if filepath.Ext(file) == ".m3u" {
			// accepted extensions
			contents = append(contents, file)
		}
	}
	sort.Strings(contents)
	return
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
