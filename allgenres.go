package main

import (
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type AllGenres []Genre

func (a *AllGenres) String() (text string) {
	text = "All Genres: \n"
	for _, genre := range *a {
		text += "\t" + genre.String()
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
		newGenre.Artists = m[genre]
		*a = append(*a, newGenre)
	}
	return
}

func (a *AllGenres) Write(path string) (err error) {
	sort.Sort(*a)
	m := make(map[string][]string)
	for _, genre := range *a {
		sort.Strings(genre.Artists)
		m[genre.Name] = genre.Artists
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, d, 0777)
	return
}
