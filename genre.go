package main

import (
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
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

func (a *AllGenres) Load(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	m := make(map[string][]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}

	for genre := range m {
		var newGenre Genre
		newGenre.Name = genre
		newGenre.Folders = m[genre]
		*a = append(*a, newGenre)
	}
	return
}

func (a *AllGenres) Write(path string) (err error) {
	sort.Sort(*a)
	m := make(map[string][]string)
	for _, genre := range *a {
		sort.Strings(genre.Folders)
		m[genre.Name] = genre.Folders
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, d, 0777)
	return
}
