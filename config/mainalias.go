package config

import (
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// MainAlias is a list of all Artists with aliases.
type MainAlias []Artist

func (a *MainAlias) String() (text string) {
	text = "All Aliases: \n"
	for _, alias := range *a {
		text += "\t" + alias.String()
	}
	return
}

func (a MainAlias) Len() int {
	return len(a)
}

func (a MainAlias) Less(i, j int) bool {
	return strings.ToLower(a[i].MainAlias) < strings.ToLower(a[j].MainAlias)
}

func (a MainAlias) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Load the configuration file where the aliases are defined.
func (a *MainAlias) Load(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	m := make(map[string][]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}

	for alias := range m {
		var newAlias Artist
		newAlias.MainAlias = alias
		newAlias.Aliases = m[alias]
		*a = append(*a, newAlias)
	}
	return
}

func (a *MainAlias) Write(path string) (err error) {
	sort.Sort(*a)
	m := make(map[string][]string)
	for _, alias := range *a {
		sort.Strings(alias.Aliases)
		m[alias.MainAlias] = alias.Aliases
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, d, 0777)
	return
}
