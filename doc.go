/*
Radis is a tool to keep your music collection in great shape.


Prerequisites

Radis works, provided your music collection is organized like this:

	root/Genre/Artist/Artist (year) Album

And that the necessary configuration files exist.


What it does

Radis can sort albums according to user-defined genres and user-defined artist
aliases.

It can track newly imported albums and adds them automatically to MPD
playlists.

These playlists can be updated if the albums move later.

It can list albums not encoded in flac, as they should all be.


More information

See http://github.com/barsanuphe/radis for more information, including how to
create the necessary configuration files.


Usage

When in doubt, try

	radis help

which should display something  like:

			R A D I S
			---------

	NAME:
	   R A D I S - Organize your music collection.

	USAGE:
	   radis [global options] command [command options] [arguments...]

	VERSION:
	   0.0.1

	COMMANDS:
	   config, c                    options for configuration
	   playlist, p                  options for playlist
	   sync, s                      sync folder according to configuration
	   check, find_awfulness        check every album is a flac version, list the heretics.
	   help, h                      Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --help, -h           show help
	   --version, -v        print the version
*/
package main
