package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/barsanuphe/radis/config"
	"github.com/barsanuphe/radis/helpers"
	"github.com/barsanuphe/radis/radis"
	"github.com/codegangsta/cli"
)

func main() {
	fmt.Printf("\n\tR A D I S\n\t---------\n\n")

	// load config
	rc := config.Config{}
	if err := rc.Load(); err != nil {
		panic(err)
	}
	// check config
	if err := rc.Check(); err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "R A D I S"
	app.Usage = "Organize your music collection."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "options for configuration",
			Subcommands: []cli.Command{
				{
					Name:    "show",
					Aliases: []string{"ls"},
					Usage:   "show configuration",
					Action: func(c *cli.Context) {
						// print config
						fmt.Println(rc.String())
					},
				},
				{
					Name:    "save",
					Aliases: []string{"sa"},
					Usage:   "reorder and save configuration files",
					Action: func(c *cli.Context) {

						if err := rc.Write(); err != nil {
							panic(err)
						}
						fmt.Println("Configuration files saved.")
					},
				},
			},
		},
		{
			Name:    "playlist",
			Aliases: []string{"p"},
			Usage:   "options for playlist",
			Subcommands: []cli.Command{
				{
					Name:    "show",
					Aliases: []string{"ls"},
					Usage:   "list playlists",
					Action: func(c *cli.Context) {
						fmt.Println("Playlists: ")
						files, err := helpers.GetPlaylistFiles(rc.Paths.MPDPlaylistDirectory)
						if err != nil {
							fmt.Println(err.Error())
						}
						for _, file := range files {
							fmt.Println(" - " + file)
						}
					},
				},
				{
					Name:    "update",
					Aliases: []string{"up"},
					Usage:   "update playlist according to configuration.",
					Action: func(c *cli.Context) {
						fmt.Println("Updating " + c.Args().First())
						p := radis.Playlist{Filename: filepath.Join(rc.Paths.MPDPlaylistDirectory, c.Args().First())}
						if err := p.UpdateAndSave(rc); err != nil {
							fmt.Println(err.Error())
						}
					},
				},
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "sync folder according to configuration",
			Action: func(c *cli.Context) {
				// sort albums
				if err := radis.SortAlbums(rc); err != nil {
					panic(err)
				}
				// scan again to remove empty directories
				if err := radis.DeleteEmptyFolders(rc); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:    "check",
			Aliases: []string{"find_awfulness"},
			Usage:   "check every album is a flac version, list the heretics.",
			Action: func(c *cli.Context) {
				// list non Flac albums
				if err := radis.FindNonFlacAlbums(rc); err != nil {
					panic(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
