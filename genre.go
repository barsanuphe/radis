package main

import "sort"

type Genre struct {
	Name    string
	Folders []string
}

func (g *Genre) String() string {
	txt := g.Name + ":\n"
	for _, artist := range g.Folders {
		txt += "\t\t- " + artist + "\n"
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
