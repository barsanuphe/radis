// Package config helps manage configuration files for radis.
package config

import (
	"fmt"

	"launchpad.net/go-xdg"
)

// Config holds the configuration for radis.
type Config struct {
	Paths   Paths
	Aliases Aliases
	Genres  Genres
}

func (c *Config) String() string {
	return c.Paths.String() + c.Aliases.String() + c.Genres.String()
}

// Check the configuration for errors.
// For now, checks the paths in radis.yaml
func (c *Config) Check() error {
	return c.Paths.Check()
}

// TODO GetConfigPath, Load, Write de radis.go
const (
	radis                  = "radis"
	radisGenresConfigFile  = radis + "_genres.yaml"
	radisAliasesConfigFile = radis + "_aliases.yaml"
	xdgMainPath            = radis + "/" + radis + ".yaml"
	xdgGenrePath           = radis + "/" + radisGenresConfigFile
	xdgAliasPath           = radis + "/" + radisAliasesConfigFile
)

func (c *Config) getConfigPaths() (mainConfigFile string, genresConfigFile string, aliasesConfigFile string, err error) {
	genresConfigFile, err = xdg.Config.Find(xdgGenrePath)
	if err != nil {
		genresConfigFile, err = xdg.Config.Ensure(xdgGenrePath)
		if err != nil {
			return
		}
		fmt.Println("Configuration file", genresConfigFile, "created. Populate it.")
	}

	aliasesConfigFile, err = xdg.Config.Find(xdgAliasPath)
	if err != nil {
		aliasesConfigFile, err = xdg.Config.Ensure(xdgAliasPath)
		if err != nil {
			return
		}
		fmt.Println("Configuration file", aliasesConfigFile, "created. Populate it.")
	}

	mainConfigFile, err = xdg.Config.Find(xdgMainPath)
	if err != nil {
		mainConfigFile, err = xdg.Config.Ensure(xdgMainPath)
		if err != nil {
			return
		}
		fmt.Println("Configuration file", mainConfigFile, "created. Populate it.")
	}
	return
}

// Load loads the configuration files into one structure, Config.
func (c *Config) Load() (err error) {
	// find configuration files
	mainConfigFile, genresConfigFile, aliasesConfigFile, err := c.getConfigPaths()
	if err != nil {
		return
	}
	// load config files
	if err = c.Paths.Load(mainConfigFile); err != nil {
		return
	}
	if err = c.Aliases.Load(aliasesConfigFile); err != nil {
		return
	}
	if err = c.Genres.Load(genresConfigFile); err != nil {
		return
	}
	return
}

// Write writes the configuration files back, after having ordered their contents.
func (c *Config) Write() (err error) {
	// find configuration files
	_, genresConfigFile, aliasesConfigFile, err := c.getConfigPaths()
	if err != nil {
		return
	}
	// write ordered config files
	if err = c.Aliases.Write(aliasesConfigFile); err != nil {
		return
	}
	if err = c.Genres.Write(genresConfigFile); err != nil {
		return
	}
	return
}
