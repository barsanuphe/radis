package config

import "testing"

var testGenres = []struct {
	genre               Genre
	expectedString      string
	candidate           string
	expectedArtist      bool
	expectedCompilation bool
}{
	// Note: Artists must be sorted for searches to work.
	// The true configuration is sorted after loading the yaml files.
	{
		Genre{Name: "test", Artists: []string{"8hopé$-ï!", "PAïof", "pif"}},
		"test:\n\t\t- 8hopé$-ï!\n\t\t- PAïof\n\t\t- pif\n",
		"8hopé$-ï!",
		true,
		false,
	},
	{
		Genre{Name: "test", Artists: []string{"hop"}},
		"test:\n\t\t- hop\n",
		"hopp",
		false,
		false,
	},
	{
		Genre{Name: "test", Artists: []string{"Hop", "Various Artists | hopp"}},
		"test:\n\t\t- Hop\n\t\t- Various Artists | hopp\n",
		"hopp",
		false,
		true,
	},
}

func TestGenreString(t *testing.T) {
	for _, ta := range testGenres {
		if v := ta.genre.String(); v != ta.expectedString {
			t.Errorf("String(%s) returned %s, expected %s!", ta.genre.Name, v, ta.expectedString)
		}
	}
}

func TestHasArtist(t *testing.T) {
	for _, ta := range testGenres {
		if v := ta.genre.HasArtist(ta.candidate); v != ta.expectedArtist {
			t.Errorf("HasArtist(%s) returned %v, expected %v!", ta.candidate, v, ta.expectedArtist)
		}
	}
}
func TestHasCompilation(t *testing.T) {
	for _, ta := range testGenres {
		if v := ta.genre.HasCompilation(ta.candidate); v != ta.expectedCompilation {
			t.Errorf("HasCompilation(%s) returned %v, expected %v!", ta.candidate, v, ta.expectedCompilation)
		}
	}
}
