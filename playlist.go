package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
		panic(err)
	}
	defer f.Close()

	// append contents
	for _, af := range p.Contents {
		// TODO: check if not already in playlist!!
		files, err := GetMusicFiles(filepath.Join(af.Root, af.NewPath))
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if _, err = f.WriteString(file + "\n"); err != nil {
				panic(err)
			}
		}
	}

	return
}

// Update a playlist by parsing the AlbumFolders it contains and writing their new paths
func (p *Playlist) Update() (err error) {
	// TODO
	return
}

// Load an existing playlist
func (p *Playlist) Load(filename string) (err error) {
	// TODO
	return
}
