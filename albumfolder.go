package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var reAlbum = regexp.MustCompile(`^([\pL\pP\pS\pN\d_ ]+) \(([0-9]+)\) ([\pL\pP\pS\pN\d_ ]+)(\s\[MP3\])?$`)

type AlbumFolder struct {
	Root      string
	Path      string
	Artist    string
	MainAlias string
	Year      string
	Title     string
	IsMP3     bool
}

func (a *AlbumFolder) String() (albumName string) {
	albumName = a.MainAlias + " / " + a.Artist + " (" + a.Year + ") " + a.Title
	if a.IsMP3 {
		albumName += " [MP3]"
	}
	return
}

func (a *AlbumFolder) IsAlbum() bool {
	if a.Artist != "" {
		// directory name already parsed, no need to do it again
		return true
	}
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
		a.MainAlias = a.Artist
		a.Year = matches[2]
		a.Title = matches[3]
		a.IsMP3 = matches[4] != ""
	} else {
		err = errors.New("Not an album!")
	}
	return
}

func (a *AlbumFolder) MoveToNewPath(genre string) (hasMoved bool, err error) {
	hasMoved = false

	if !a.IsAlbum() {
		return false, errors.New("Cannot move, not an album.")
	}

	directoryName := filepath.Base(a.Path)
	newPath := filepath.Join(a.Root, genre, a.MainAlias, directoryName)
	// comparer avec l'ancien
	if newPath != a.Path {
		// if different, move folder
		originalRelative, _ := filepath.Rel(a.Root, a.Path)
		destRelative, _ := filepath.Rel(a.Root, newPath)
		fmt.Println("- "+originalRelative, " -> ", destRelative)

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
		if err == nil {
			hasMoved = true
		}
	}
	return
}
