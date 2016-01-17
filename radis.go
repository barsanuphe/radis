package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"launchpad.net/go-xdg"
)

// CONFIG ----------------------------------------------------------------------

const (
	radis                  = "radis"
	radisGenresConfigFile  = radis + "_genres.yaml"
	radisAliasesConfigFile = radis + "_aliases.yaml"
	xdgGenrePath           = radis + "/" + radisGenresConfigFile
	xdgAliasPath           = radis + "/" + radisAliasesConfigFile
)

func getConfigPaths() (genresConfigFile string, aliasesConfigFile string, err error) {
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
	return
}

func LoadConfig() (aliases MainAlias, genres AllGenres, err error) {
	// find configuration files
	genresConfigFile, aliasesConfigFile, err := getConfigPaths()
	if err != nil {
		return
	}
	// load config files
	aliases = MainAlias{}
	if err = aliases.Load(aliasesConfigFile); err != nil {
		return
	}
	genres = AllGenres{}
	if err = genres.Load(genresConfigFile); err != nil {
		return
	}
	return
}

func WriteConfig(aliases MainAlias, genres AllGenres) (err error) {
	// find configuration files
	genresConfigFile, aliasesConfigFile, err := getConfigPaths()
	if err != nil {
		return
	}
	// write ordered config files
	if err = aliases.Write(aliasesConfigFile); err != nil {
		return
	}
	if err = genres.Write(genresConfigFile); err != nil {
		return
	}
	return
}

// HELPERS ---------------------------------------------------------------------

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("-- [%s done in %s]\n", name, elapsed)
}

func GetExistingPath(path string) (existingPath string, err error) {
	// check root exists or pwd+root exists
	if filepath.IsAbs(path) {
		existingPath = path
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		existingPath = filepath.Join(pwd, path)
	}
	// check root exists
	if _, err = os.Stat(existingPath); os.IsNotExist(err) {
		err = errors.New("Directory " + path + " does not exist!!!")
	}
	return
}

// SORT ------------------------------------------------------------------------

func sortAlbums(root string, aliases MainAlias, genres AllGenres) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for albums in " + root + ".")
	movedAlbums := 0
	uncategorized := 0
	foundAlbums := 0
	mp3Albums := 0
	err = filepath.Walk(root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: root, Path: path}
			if af.IsAlbum() {
				foundAlbums++
				if af.IsMP3 {
					mp3Albums++
				}
				hasMoved := false
				// fmt.Println("+ Found album: ", af.String())
				found := false

				// see if artist has known alias
				for _, alias := range aliases {
					if alias.HasAlias(af.Artist) {
						af.MainAlias = alias.MainAlias
						break
					}
				}
				// find which genre the artist or main alias belongs to
				for _, genre := range genres {
					// if artist is known, it belongs to genre.Name
					if genre.HasArtist(af.MainAlias) || genre.HasCompilation(af.Title) {
						hasMoved, err = af.MoveToNewPath(genre.Name)
						found = true
						break
					}
				}
				if !found {
					uncategorized++
					hasMoved, err = af.MoveToNewPath("UNCATEGORIZED")
				}
				if hasMoved {
					movedAlbums++
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
	return
}

// LIST NON FLAC ---------------------------------------------------------------

func HasNonFlacFiles(directoryPath string) (bool, error) {
	f, err := os.Open(directoryPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fileList, err := f.Readdirnames(-1)
	if err == io.EOF {
		return true, nil
	}

	hasNonFlac := false
	for _, file := range fileList {
		switch filepath.Ext(file) {
		case ".flac", ".jpg", ".jpeg", ".png":
			// accepted extensions
		case ".mp3", ".wma", ".m4a":
			hasNonFlac = true
			break
		default:
			fmt.Println("Found suspicious file ", file, " in ", directoryPath)
			hasNonFlac = true
			break
		}
	}
	return hasNonFlac, err
}

func findNonFlacAlbums(root string) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for non-Flac albums in " + root + ".")
	unFlagged := 0
	nonFlacAlbums := 0
	err = filepath.Walk(root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
		// when an album has just been moved, Walk goes through it a second
		// time with an "file does not exist" error
		if os.IsNotExist(walkError) {
			return
		}

		if fileInfo.IsDir() {
			af := AlbumFolder{Root: root, Path: path}
			if af.IsAlbum() {
				// scan contents for non-flac
				isNonFlac, err := HasNonFlacFiles(path)
				if err != nil {
					panic(err)
				}
				relativePath, _ := filepath.Rel(root, path)
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

// CLEAN -----------------------------------------------------------------------

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func deleteEmptyFolders(root string) (err error) {
	defer timeTrack(time.Now(), "Scanning files")

	fmt.Println("Scanning for empty directories.")
	deletedDirectories := 0
	deletedDirectoriesThisTime := 0
	atLeastOnce := false

	// loops until all levels of empty directories are deleted
	for !atLeastOnce || deletedDirectoriesThisTime != 0 {
		atLeastOnce = true
		deletedDirectoriesThisTime = 0
		err = filepath.Walk(root, func(path string, fileInfo os.FileInfo, walkError error) (err error) {
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

// MAIN ------------------------------------------------------------------------

func main() {
	fmt.Println("\n\tR A D I S\n\t---------\n")

	// load config
	aliases, genres, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	// cli: commands show / sync folder
	app := cli.NewApp()
	app.Name = "R A D I S"
	app.Usage = "Organize your music collection."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"ls"},
			Usage:   "show configuration",
			Action: func(c *cli.Context) {
				// print config
				fmt.Println(aliases.String())
				fmt.Println(genres.String())
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "sync folder according to configuration",
			Action: func(c *cli.Context) {
				// scan folder in root
				root, err := GetExistingPath(c.Args().First())
				if err != nil {
					panic(err)
				}

				// sort albums
				if err := sortAlbums(root, aliases, genres); err != nil {
					panic(err)
				}
				// scan again to remove empty directories
				if err := deleteEmptyFolders(root); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:    "check",
			Aliases: []string{"find_awfulness"},
			Usage:   "check every album is a flac version, list the heretics.",
			Action: func(c *cli.Context) {
				// scan folder in root
				root, err := GetExistingPath(c.Args().First())
				if err != nil {
					panic(err)
				}

				// list non Flac albums
				if err := findNonFlacAlbums(root); err != nil {
					panic(err)
				}
			},
		},
	}

	app.Run(os.Args)

	// write ordered config files
	if err := WriteConfig(aliases, genres); err != nil {
		panic(err)
	}
}
