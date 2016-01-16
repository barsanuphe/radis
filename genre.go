package main

import "sort"

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

func (g *Genre) HasArtist(artist string) bool {
	sort.Strings(g.Artists)
	i := sort.SearchStrings(g.Artists, artist)
	if i < len(g.Artists) && g.Artists[i] == artist {
		// fmt.Println("++ Found artist ", artist, "in genre ", g.Name)
		return true
	}
	return false
}
