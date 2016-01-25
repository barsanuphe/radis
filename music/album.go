// Package music deals with album folders and m3u playlists.
package music

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/barsanuphe/radis/config"
	"github.com/barsanuphe/radis/directory"

)

var albumPattern = regexp.MustCompile(`^([\pL\pP\pS\pN\d\pZ]+) \(([0-9]{4})\) ([\pL\pP\pS\pN\d\pZ]+?)(\[MP3\])?$`)

// Album holds the information of an album directory.
// An album follows the pattern: Artist (year) Album title
// Or: Various Artists (year) Compilation title
type Album struct {
	Root      string // absolute
	Path      string // absolute
	NewPath   string // absolute
	artist    string
	mainAlias string
	year      string
	title     string
	IsMP3     bool
}

// String gives a representation of an AlbumFolder.
func (a *Album) String() (albumName string) {
	albumName = a.mainAlias + "/" + a.artist + " (" + a.year + ") " + a.title
	if a.IsMP3 {
		albumName += "[MP3]"
	}
	return
}

// IsValidAlbum indicates if a directory name has the proper template to be an album.
func (a *Album) IsValidAlbum() bool {
	if a.artist != "" {
		// directory name already parsed, no need to do it again
		return true
	}
	if err := a.extractInfo(); err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

// IsNew checks if the album was found in the INCOMING directory.
func (a *Album) IsNew(c config.Config) bool {
	return strings.Contains(a.Path, filepath.Join(c.Paths.Root, c.Paths.IncomingSubdir))
}

// extractInfo parses an AlbumFolder's basepath to extract information.
func (a *Album) extractInfo() (err error) {
	matches := albumPattern.FindStringSubmatch(filepath.Base(a.Path))
	if len(matches) > 0 {
		a.artist = matches[1]
		a.mainAlias = a.artist
		a.year = matches[2]
		a.title = matches[3]
		a.IsMP3 = matches[4] != ""
	} else {
		err = errors.New("Not an album!")
	}
	return
}

// FindNewPath for an album according to configuration.
func (a *Album) FindNewPath(c config.Config) (hasGenre bool, err error) {
	if !a.IsValidAlbum() {
		err = errors.New("Not an album!")
		return
	}

	// see if artist has known alias
	for _, alias := range c.Aliases {
		if alias.HasAlias(a.artist) {
			a.mainAlias = alias.MainAlias
			break
		}
	}
	// find which genre the artist or main alias belongs to
	hasGenre = false
	directoryName := filepath.Base(a.Path)
	for _, genre := range c.Genres {
		var found bool
		if a.mainAlias == "Various Artists" {
			found = genre.HasCompilation(a.title)
		} else {
			found = genre.HasArtist(a.mainAlias)
		}
		// if artist is known, it belongs to genre.Name
		if found {
			a.NewPath = filepath.Join(a.Root, genre.Name, a.mainAlias, directoryName)
			hasGenre = true
			break
		}
	}
	if !hasGenre {
		a.NewPath = filepath.Join(a.Root, c.Paths.UnsortedSubdir, a.mainAlias, directoryName)
	}
	return
}

// MoveToNewPath moves an album directory to its new home in another genre.
func (a *Album) MoveToNewPath(doNothing bool) (hasMoved bool, err error) {
	hasMoved = false
	if a.NewPath == "" {
		return false, errors.New("FindNewPath first.")
	}
	// comparer avec l'ancien
	if a.NewPath != a.Path {
		// if different, move folder
		if !doNothing {
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
		} else {
			// would have moved, but must do nothing
			hasMoved = true
		}
	}
	return
}

// GetMusicFiles returns flac or mp3 files of the album.
func (a *Album) GetMusicFiles() (contents []string, err error) {
	fileList, err := directory.GetFiles(a.NewPath)
	if err != nil {
		return []string{}, err
	}
	// check for music files
	for _, file := range fileList {
		switch filepath.Ext(file) {
		case ".flac", ".mp3":
			// accepted extensions
			contents = append(contents, filepath.Join(a.NewPath, file))
		}
	}
	sort.Strings(contents)
	return
}

// HasNonFlacFiles returns true if an album contains files other than flac songs and cover pictures.
func (a *Album) HasNonFlacFiles() (bool, error) {
	fileList, err := directory.GetFiles(a.Path)
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
			fmt.Println("Found suspicious file ", file, " in ", a.Path)
			hasNonFlac = true
			break
		}
	}
	return hasNonFlac, err
}
