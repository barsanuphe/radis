// Radis is a tool to keep your music collection in great shape.
//
// That is, provided your music collection is organized like this:
//
//	root/Genre/Artist/Artist (year) Album
//
// see github.com/barsanuphe/radis for more information.
//
// Usage:
//
// This command lists what was found in the configuration files:
//
//    $ radis show
//
// This reorganizes your music collection in the "Root" indicated in radis.yaml:
//
//   $ radis sync
//
// Make sure "Root" is correct.
// radis will stop if the path does not exist, but otherwise it will at least
// delete empty directories in that "Root".
//
// Of course, you should only have flac versions of your music.
// Sometimes they do not exist, so these albums have a "[MP3]" suffix in the
// folder name.
//
// To list those offending albums, and check you have not missed any:
//    $ radis check
//
// When in doubt:
//
//   $ radis help
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barsanuphe/radis/config"
	"github.com/codegangsta/cli"
)

// sortAlbums scans the music collection root and reorders albums according to the configuration files.
func sortAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	movedAlbums := 0
	uncategorized := 0
	foundAlbums := 0
	mp3Albums := 0

	dailyPlaylist, monthlyPlaylist := loadCurrentPlaylists(c)

	fmt.Println("Scanning for albums in " + c.Paths.Root + ".")
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: c.Paths.Root, Path: path}
			if af.IsAlbum() {
				foundAlbums++
				if af.IsMP3 {
					mp3Albums++
				}
				hasGenre, err := af.FindNewPath(c)
				if err != nil {
					panic(err)
				}
				if !hasGenre {
					uncategorized++
				}
				hasMoved, err := af.MoveToNewPath()
				if err != nil {
					panic(err)
				}
				if hasMoved {
					movedAlbums++
				}

				// find out if we should append playlists
				isInsideIncoming := strings.Contains(path, filepath.Join(c.Paths.Root, c.Paths.IncomingSubdir))
				if isInsideIncoming {
					// album is inside INCOMING dir, add to playlist automatically
					fmt.Println("  ++ " + af.String() + " was in INCOMING, adding to playlist.")
					dailyPlaylist.Contents = append(dailyPlaylist.Contents, af)
					monthlyPlaylist.Contents = append(monthlyPlaylist.Contents, af)
				}
			}
		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Printf("Found %d albums (%d MP3 albums), Moved %d.\n", foundAlbums, mp3Albums, movedAlbums)
	if uncategorized != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", uncategorized)
	}

	if err := writeCurrentPlaylists(dailyPlaylist, monthlyPlaylist); err != nil {
		panic(err)
	}
	return
}

// findNonFlacAlbums scan the music collection root and lists all albums of mp3 files instead of flac.
func findNonFlacAlbums(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for non-Flac albums in " + c.Paths.Root + ".")
	unFlagged := 0
	nonFlacAlbums := 0
	err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: c.Paths.Root, Path: path}
			if af.IsAlbum() {
				// scan contents for non-flac
				isNonFlac, err := HasNonFlacFiles(path)
				if err != nil {
					panic(err)
				}
				relativePath, _ := filepath.Rel(c.Paths.Root, path)
				if isNonFlac {
					fmt.Println("- ", relativePath)
					nonFlacAlbums++
				}
				if isNonFlac && !af.IsMP3 {
					unFlagged++
					fmt.Println("!!! ", relativePath, " not flagged as non FLAC!!!")
				}
				// NOTE: find falsely tagged folders? is that a thing?
			}
		}
		return
	})
	if err != nil {
		fmt.Printf("Error!")
	}
	fmt.Printf("Found %d non-Flac albums, including %d incorrectly flagged.\n", nonFlacAlbums, unFlagged)
	if unFlagged != 0 {
		fmt.Printf("\n!!!\n!!! %d album(s) remain UNCATEGORIZED !!!\n!!!\n\n", unFlagged)
	}
	return
}

// deleteEmptyFolders deletes empty folders that may appear after sorting albums.
func deleteEmptyFolders(c config.Config) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for empty directories.")
	deletedDirectories := 0
	deletedDirectoriesThisTime := 0
	atLeastOnce := false

	// loops until all levels of empty directories are deleted
	for !atLeastOnce || deletedDirectoriesThisTime != 0 {
		atLeastOnce = true
		deletedDirectoriesThisTime = 0
		err = filepath.Walk(c.Paths.Root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
			// when an album has just been removed, Walk goes through it a second
			// time with an "file does not exist" error
			if os.IsNotExist(walkError) {
				return
			}
			if fileInfo.IsDir() {
				isEmpty, err := IsEmpty(path)
				if err != nil {
					panic(err)
				}
				if isEmpty {
					fmt.Println("Removing empty directory ", path)
					if err := os.Remove(path); err == nil {
						deletedDirectories++
						deletedDirectoriesThisTime++
					}
				}
			}
			return
		})
		if err != nil {
			fmt.Printf("Error!")
		}
	}

	fmt.Printf("Removed %d albums.\n", deletedDirectories)
	return
}

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
						files, err := GetPlaylistFiles(rc.Paths.MPDPlaylistDirectory)
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
						p := Playlist{Filename: filepath.Join(rc.Paths.MPDPlaylistDirectory, c.Args().First())}
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
				if err := sortAlbums(rc); err != nil {
					panic(err)
				}
				// scan again to remove empty directories
				if err := deleteEmptyFolders(rc); err != nil {
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
				if err := findNonFlacAlbums(rc); err != nil {
					panic(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
