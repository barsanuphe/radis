package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	f, err := os.OpenFile(p.Filename, os.O_WRONLY|os.O_CREATE, 0600)
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
		// nothing to do
		return
	}
	for i := range p.Contents {
		if err := p.Contents[i].ExtractInfo(); err != nil {
			panic(err)
		}
		// find the new path, so that it can be exported by Write
		if _, err := p.Contents[i].FindNewPath(c); err != nil {
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
func (p *Playlist) Load(root string) (err error) {
	// open file and get strings
	content, err := ioutil.ReadFile(p.Filename)
	if os.IsNotExist(err) {
		// file does not exist, nothing to do
		return nil
	} else if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	//remove filename
	albumsPaths := []string{}
	for _, l := range lines {
		if l != "" {
			albumsPaths = append(albumsPaths, filepath.Dir(l))
		}
	}
	// remove duplicates
	albumsPaths = removeDuplicatePaths(albumsPaths)

	// add AlbumFolder to Contents
	for _, a := range albumsPaths {
		p.Contents = append(p.Contents, AlbumFolder{Root: root, Path: a})
	}
	return
}

// loadCurrentPlaylists finds and loads current playlists
func loadCurrentPlaylists(c config.Config) (daily Playlist, monthly Playlist) {
	now := time.Now().Local()
	thisDay := now.Format("2006-01-02")
	thisMonth := now.Format("2006-01")

	daily = Playlist{
		Filename: filepath.Join(c.Paths.MPDPlaylistDirectory, thisDay+".m3u"),
	}
	monthly = Playlist{
		Filename: filepath.Join(c.Paths.MPDPlaylistDirectory, thisMonth+".m3u"),
	}

	// Load the playlists if they exist
	err := daily.Load(c.Paths.Root)
	if err != nil {
		panic(err)
	}
	err = monthly.Load(c.Paths.Root)
	if err != nil {
		panic(err)
	}

	// Update the playlists if they exist
	err = daily.Update(c)
	if err != nil {
		panic(err)
	}
	err = monthly.Update(c)
	if err != nil {
		panic(err)
	}
	return
}

// writePlaylists after sync
func writePlaylists(daily Playlist, monthly Playlist) (err error) {
	if len(daily.Contents) != 0 {
		fmt.Println("Writing playlist " + filepath.Base(daily.Filename) + ".")
		err = daily.Write()
		if err != nil {
			return
		}
		fmt.Println("Writing playlist " + filepath.Base(monthly.Filename) + ".")
		err = monthly.Write()
		if err != nil {
			return
		}
	}
	return
}
