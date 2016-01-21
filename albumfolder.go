package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/barsanuphe/radis/config"
)

var albumPattern = regexp.MustCompile(`^([\pL\pP\pS\pN\d\pZ]+) \(([0-9]{4})\) ([\pL\pP\pS\pN\d\pZ]+?)(\[MP3\])?$`)

// AlbumFolder holds the information of an album directory.
// An album follows the pattern: Artist (year) Album title
type AlbumFolder struct {
	Root      string // absolute
	Path      string // absolute
	NewPath   string // absolute
	Artist    string
	MainAlias string
	Year      string
	Title     string
	IsMP3     bool
}

// String gives a representation of an AlbumFolder.
func (a *AlbumFolder) String() (albumName string) {
	albumName = a.MainAlias + "/" + a.Artist + " (" + a.Year + ") " + a.Title
	if a.IsMP3 {
		albumName += "[MP3]"
	}
	return
}

// IsAlbum indicates if a directory name has the proper template to be an album.
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

// ExtractInfo parses an AlbumFolder's basepath to extract information.
func (a *AlbumFolder) ExtractInfo() (err error) {
	matches := albumPattern.FindStringSubmatch(filepath.Base(a.Path))
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

// FindNewPath for an album according to configuration.
func (a *AlbumFolder) FindNewPath(c config.Config) (hasGenre bool, err error) {
	if !a.IsAlbum() {
		err = errors.New("Not an album!")
		return
	}

	// see if artist has known alias
	for _, alias := range c.Aliases {
		if alias.HasAlias(a.Artist) {
			a.MainAlias = alias.MainAlias
			break
		}
	}
	// find which genre the artist or main alias belongs to
	hasGenre = false
	directoryName := filepath.Base(a.Path)
	for _, genre := range c.Genres {
		// if artist is known, it belongs to genre.Name
		if genre.HasArtist(a.MainAlias) || genre.HasCompilation(a.Title) {
			a.NewPath = filepath.Join(a.Root, genre.Name, a.MainAlias, directoryName)
			hasGenre = true
			break
		}
	}
	if !hasGenre {
		a.NewPath = filepath.Join(a.Root, c.Paths.UnsortedSubdir, a.MainAlias, directoryName)
	}
	return
}

// MoveToNewPath moves an album directory to its new home in another genre.
func (a *AlbumFolder) MoveToNewPath() (hasMoved bool, err error) {
	hasMoved = false
	if a.NewPath == "" {
		return false, errors.New("FindNewPath first.")
	}
	// comparer avec l'ancien
	if a.NewPath != a.Path {
		// if different, move folder
		originalRelative, _ := filepath.Rel(a.Root, a.Path)
		destRelative, _ := filepath.Rel(a.Root, a.NewPath)
		fmt.Println("- "+originalRelative, " -> ", destRelative)

		newPathParent := filepath.Dir(a.NewPath)
		if _, err = os.Stat(newPathParent); os.IsNotExist(err) {
			// newPathParent does not exist, creating
			err = os.MkdirAll(newPathParent, 0777)
			if err != nil {
				panic(err)
			}
		}
		// move
		err = os.Rename(a.Path, a.NewPath)
		if err == nil {
			hasMoved = true
		}
	}
	return
}
