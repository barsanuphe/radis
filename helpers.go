package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"sort"
)

// timeTrack can be used to evaluate the time spent in a function.
// usage: defer timeTrack(startTime) at the beginning of the function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
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
		err = errors.New("Directory " + path + " does not exist!!!")
	}
	return
}

// HasNonFlacFiles returns true if an album contains files other than flac songs and cover pictures.
func HasNonFlacFiles(albumPath string) (bool, error) {
	f, err := os.Open(albumPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// get contents
	fileList, err := f.Readdirnames(-1)
	if err == io.EOF {
		return true, nil
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
	f, err := os.Open(albumPath)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	// get contents
	fileList, err := f.Readdirnames(-1)
	if err == io.EOF {
		return []string{}, nil
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
