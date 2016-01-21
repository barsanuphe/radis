package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/barsanuphe/radis/config"
)

// Playlist can generate .m3u playlists from a list of AlbumFolders.
type Playlist struct {
	Filename string
	Contents []AlbumFolder
}

// String gives a representation of an Playlist.
func (p *Playlist) String() (playlist string) {
	playlist = fmt.Sprintf("%s: %d albums", p.Filename, len(p.Contents))
	return
}

// Write the playlist or append it if it exists
func (p *Playlist) Write() (err error) {
	if len(p.Contents) == 0 {
		err = errors.New("Empty playlist, nothing to write.")
		return
	}
	// open
	f, err := os.OpenFile(p.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	// remove duplicates
	if err := p.RemoveDuplicates(); err != nil {
		return err
	}

	// append contents
	for _, af := range p.Contents {
		// TODO: in radis, load existing playlists, add new albums
		// TODO: before write, Update()

		files, err := GetMusicFiles(af.NewPath)
		if err != nil {
			return err
		}
		for _, file := range files {
			if _, err = f.WriteString(file + "\n"); err != nil {
				return err
			}
		}
	}
	return
}

// Update a playlist by parsing the AlbumFolders it contains and writing their new paths
func (p *Playlist) Update(c config.Config) (err error) {
	if len(p.Contents) == 0 {
		return errors.New("Empty playlist.")
	}
	for _, af := range p.Contents {
		if err := af.ExtractInfo(); err != nil {
			panic(err)
		}
		// find the new path, so that it can be exported by Write
		if _, err := af.FindNewPath(c); err != nil {
			panic(err)
		}
	}
	return
}

// RemoveDuplicates in a Playlist's Contents
func (p *Playlist) RemoveDuplicates() (err error) {
	result := []AlbumFolder{}
	seen := map[string]AlbumFolder{}
	for _, val := range p.Contents {
		if _, ok := seen[val.Path]; !ok {
			result = append(result, val)
			seen[val.Path] = val
		}
	}
	// replace contents
	p.Contents = result
	return
}

// Load an existing playlist
func (p *Playlist) Load(filename string, root string) (err error) {
	// open file and get strings
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	//remove filename
	albumsPaths := []string{}
	for _, l := range lines {
		albumsPaths = append(albumsPaths, filepath.Dir(l))
	}
	// remove duplicates
	albumsPaths = removeDuplicatePaths(albumsPaths)

	// add AlbumFolder to Contents
	for _, a := range albumsPaths {
		p.Contents = append(p.Contents, AlbumFolder{Root: root, Path: a})
	}
	return
}
