package main

import "testing"

var af1 = AlbumFolder{Root: ".", Path: "test (2009) doié? [MP3]"}
var af2 = AlbumFolder{Root: ".", Path: ".._12Jïâ!! (2012) AA1AQ"}
var p = Playlist{Filename: "hop.m3u", Contents: []AlbumFolder{af1, af2}}

var testPlaylists = []struct {
	Playlist Playlist
	expected string
}{
	{p, "hop.m3u: 2 albums"},
	{Playlist{Filename: "éé?.m3u", Contents: []AlbumFolder{af1}}, "éé?.m3u: 1 albums"},
}

func TestPlaylistString(t *testing.T) {
	for _, tp := range testPlaylists {
		if st := tp.Playlist.String(); st != tp.expected {
			t.Errorf("String(%s) returned %s, expected %s", tp.Playlist.Filename, st, tp.expected)
		}
	}
}

func TestWrite(t *testing.T) {
	for _, tp := range testPlaylists {
		// TODO create bogus files so that write actually does something
		if err := tp.Playlist.Write(); err != nil && err.Error() != "Could not find path ; have you synced lately?" {
			t.Errorf("Write(%s) returned %s, expected nil", tp.Playlist.Filename, err.Error())
		}
		// TODO check file contents
		/*if err := os.Remove(tp.Playlist.Filename); err != nil {
			t.Errorf("Could not cleanup playlist!")
		}*/
	}
}

func TestUpdate(t *testing.T) {
	// TODO
}

func TestLoad(t *testing.T) {
	// TODO
}
