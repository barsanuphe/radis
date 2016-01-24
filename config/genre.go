package config

import "sort"

// Genre is a struct defining a genre and the artists that belong to it.
type Genre struct {
	Name    string
	Artists []string
}

func (g *Genre) String() string {
	txt := g.Name + ":\n"
	for _, artist := range g.Artists {
		txt += "\t\t- " + artist + "\n"
	}
	return txt
}

// HasArtist checks if the Genre contains a given artist.
func (g *Genre) HasArtist(artist string) bool {
	// already sorted at Load
	i := sort.SearchStrings(g.Artists, artist)
	if i < len(g.Artists) && g.Artists[i] == artist {
		// fmt.Println("++ Found artist ", artist, "in genre ", g.Name)
		return true
	}
	return false
}

// HasCompilation checks if the Genre contains a compilation with a specific title.
func (g *Genre) HasCompilation(title string) bool {
	fullTitle := "Various Artists | " + title
	// already sorted at Load
	i := sort.SearchStrings(g.Artists, fullTitle)
	if i < len(g.Artists) && g.Artists[i] == fullTitle {
		// fmt.Println("++ Found compilation ", fullTitle, "in genre ", g.Name)
		return true
	}
	return false
}
