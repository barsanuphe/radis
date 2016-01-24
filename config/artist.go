package config

import "sort"

// Artist can define the main alias for artists that have more than one.
type Artist struct {
	MainAlias string
	Aliases   []string
}

func (g *Artist) String() string {
	txt := g.MainAlias + ":\n"
	for _, alias := range g.Aliases {
		txt += "\t\t- " + alias + "\n"
	}
	return txt
}

// HasAlias can check if an Artist has a given alias.
func (g *Artist) HasAlias(alias string) bool {
	// already sorted at Load
	i := sort.SearchStrings(g.Aliases, alias)
	if i < len(g.Aliases) && g.Aliases[i] == alias {
		// fmt.Println("++ Found alias ", alias, "in genre ", g.MainAlias)
		return true
	}
	return false
}
