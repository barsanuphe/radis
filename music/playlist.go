package music

import (
	"errors"
	"fmt"
	"github.com/barsanuphe/radis/config"
	"github.com/barsanuphe/radis/directory"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Playlist can generate .m3u playlists from a list of AlbumFolders.
type Playlist struct {
	Filename string
	contents []AlbumFolder
}

// String gives a representation of an Playlist.
func (p *Playlist) String() (playlist string) {
	playlist = fmt.Sprintf("%s: %d albums", p.Filename, len(p.contents))
	return
}

// Exists finds out if a Playlist is valid or not
func (p *Playlist) Exists() (isPlaylist bool, err error) {
	path, err := directory.GetExistingPath(p.Filename)
	if err != nil {
		return
	}
	if filepath.Ext(path) == ".m3u" {
		isPlaylist = true
	}
	return
}

// RemoveDuplicates in a Playlist's Contents
func (p *Playlist) RemoveDuplicates() (err error) {
	result := []AlbumFolder{}
	seen := map[string]AlbumFolder{}
	for _, val := range p.contents {
		if _, ok := seen[val.Path]; !ok {
			result = append(result, val)
			seen[val.Path] = val
		}
	}
	// replace contents
	p.contents = result
	return
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
		p.contents = append(p.contents, AlbumFolder{Root: root, Path: a})
	}
	return
}

// Update a playlist by parsing the AlbumFolders it contains and writing their new paths
func (p *Playlist) Update(c config.Config) (err error) {
	if len(p.contents) == 0 {
		// nothing to do
		return
	}
	for i := range p.contents {
		if err := p.contents[i].ExtractInfo(); err != nil {
			panic(err)
		}
		// find the new path, so that it can be exported by Write
		if _, err := p.contents[i].FindNewPath(c); err != nil {
			panic(err)
		}
	}
	return
}

// Write the playlist or append it if it exists
func (p *Playlist) Write() (err error) {
	if len(p.contents) == 0 {
		err = errors.New("Empty playlist, nothing to write.")
		return
	}

	// remove duplicates
	if err := p.RemoveDuplicates(); err != nil {
		return err
	}

	// append contents
	contents := []string{}
	for _, af := range p.contents {
		files, err := af.GetMusicFiles()
		if os.IsNotExist(err) {
			return errors.New("Could not find path " + af.NewPath + "; have you synced lately?")
		} else if err != nil {
			return err
		}
		for i := range files {
			// MPD wants relative paths
			relativePath, err := filepath.Rel(af.Root, files[i])
			if err != nil {
				panic(err)
			}
			contents = append(contents, relativePath)
		}
	}

	// write if everything is good.
	f, err := os.OpenFile(p.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	for _, file := range contents {
		if _, err = f.WriteString(file + "\n"); err != nil {
			return err
		}
	}
	return
}

// UpdateAndSave a Playlist file.
func (p *Playlist) UpdateAndSave(c config.Config) (err error) {
	isPlaylist, err := p.Exists()
	if err != nil {
		return
	} else if !isPlaylist {
		return errors.New(p.Filename + "does not exist!")
	}
	// Load the playlist
	err = p.Load(c.Paths.Root)
	if err != nil {
		panic(err)
	}
	// Update the playlist
	err = p.Update(c)
	if err != nil {
		panic(err)
	}
	// Write the playlist
	err = p.Write()
	if err != nil {
		return
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
func writeCurrentPlaylists(daily Playlist, monthly Playlist) (err error) {
	if len(daily.contents) != 0 {
		fmt.Println("Writing playlist " + filepath.Base(daily.Filename) + ".")
		err = daily.Write()
		if err != nil {
			return
		}
	}
	if len(monthly.contents) != 0 {
		fmt.Println("Writing playlist " + filepath.Base(monthly.Filename) + ".")
		err = monthly.Write()
		if err != nil {
			return
		}
	}
	return
}
