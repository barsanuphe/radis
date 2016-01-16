package main

import (
	"sort"
	"strings"
)

type Genre struct {
	Name    string
	Folders []string
}

func (g *Genre) String() string {
	txt := g.Name + ":\n"
	for _, artist := range g.Folders {
		txt += "\t- " + artist + "\n"
	}
	return txt
}

func (g *Genre) HasArtist(artist string) bool {
	sort.Strings(g.Folders)
	i := sort.SearchStrings(g.Folders, artist)
	if i < len(g.Folders) && g.Folders[i] == artist {
		// fmt.Println("++ Found artist ", artist, "in genre ", g.Name)
		return true
	}
	return false
}

//--------------------------

type AllGenres []Genre

func (a *AllGenres) String() (text string) {
	text = "All Genres: \n"
	for _, genre := range *a {
		text += genre.String()
	}
	return
}

func (a AllGenres) Len() int {
	return len(a)
}

func (a AllGenres) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
}

func (a AllGenres) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
