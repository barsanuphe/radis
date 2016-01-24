// Package directory helps deal with folder structure.
package directory

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"


)

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

// GetFiles returns the files inside a path
func GetFiles(path string) (contents []string, err error) {
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

// GetPlaylistFiles returns a list of .m3u files.
func GetPlaylists(playlistRoot string) (contents []string, err error) {
	fileList, err := GetFiles(playlistRoot)
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
