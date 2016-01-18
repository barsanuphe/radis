package config

import "sort"

type Artist struct {
	MainAlias    string
	Aliases []string
}

func (g *Artist) String() string {
	txt := g.MainAlias + ":\n"
	for _, alias := range g.Aliases {
		txt += "\t\t- " + alias + "\n"
	}
	return txt
}

func (g *Artist) HasAlias(alias string) bool {
	sort.Strings(g.Aliases)
	i := sort.SearchStrings(g.Aliases, alias)
	if i < len(g.Aliases) && g.Aliases[i] == alias {
		// fmt.Println("++ Found alias ", alias, "in genre ", g.MainAlias)
		return true
	}
	return false
}
