package config

import (
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// Aliases is a list of all Artists with aliases.
type Aliases []Artist

func (a *Aliases) String() (text string) {
	text = "All Aliases: \n"
	for _, alias := range *a {
		text += "\t" + alias.String()
	}
	return
}

func (a Aliases) Len() int {
	return len(a)
}

func (a Aliases) Less(i, j int) bool {
	return strings.ToLower(a[i].MainAlias) < strings.ToLower(a[j].MainAlias)
}

func (a Aliases) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Load the configuration file where the aliases are defined.
func (a *Aliases) Load(path string) (err error) {
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
		sort.Strings(m[alias])
		newAlias.Aliases = m[alias]
		*a = append(*a, newAlias)
	}
	sort.Sort(*a)
	return
}

func (a *Aliases) Write(path string) (err error) {
	m := make(map[string][]string)
	for _, alias := range *a {
		m[alias.MainAlias] = alias.Aliases
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, d, 0777)
	return
}
