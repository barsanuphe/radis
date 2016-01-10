package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Genre struct {
	Name    string
	Folders []string
}

func (g *Genre) String() (string) {
	txt := g.Name + ":\n"
	for _, artist := range g.Folders{
		txt += "\t- " + artist + "\n"
	}
	return txt
}

//----------------

func printConfig(config []Genre) {
		for _ , genre := range config {
		fmt.Println(genre.String())
	}
}

func readConfig(path string) (Config []Genre, err error) {
	data, err := ioutil.ReadFile("radis.yaml")
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
		Config = append(Config, newGenre)
	}
	return
}


//----------------

func main() {
	fmt.Println("R A D I S\n---------\n")

	Config, err := readConfig("radis.yaml")
	if err != nil {
		panic(err)
	}

	printConfig(Config)

}
