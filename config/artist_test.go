package config

import "testing"

var testArtists = []struct {
	artist         Artist
	expectedString string
	candidate      string
	expectedAlias  bool
}{
	{
		Artist{MainAlias: "test", Aliases: []string{"8hopé$-ï!", "paf", "pif"}},
		"test:\n\t\t- 8hopé$-ï!\n\t\t- paf\n\t\t- pif\n",
		"8hopé$-ï!",
		true,
	},
	{
		Artist{MainAlias: "test", Aliases: []string{"hop"}},
		"test:\n\t\t- hop\n",
		"hopp",
		false,
	},
}

func TestArtistString(t *testing.T) {
	for _, ta := range testArtists {
		if v := ta.artist.String(); v != ta.expectedString {
			t.Errorf("String(%s) returned %s, expected %s!", ta.artist.MainAlias, v, ta.expectedString)
		}
	}
}

func TestHasAlias(t *testing.T) {
	for _, ta := range testArtists {
		if v := ta.artist.HasAlias(ta.candidate); v != ta.expectedAlias {
			t.Errorf("HasAlias(%s) returned %v, expected %v!", ta.candidate, v, ta.expectedAlias)
		}
	}
}
