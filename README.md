[![GoDoc](https://godoc.org/github.com/barsanuphe/radis?status.svg)](https://godoc.org/github.com/barsanuphe/radis)
[![Build Status](https://travis-ci.org/barsanuphe/radis.svg?branch=master)](https://travis-ci.org/barsanuphe/radis)
[![Report Card](http://goreportcard.com/badge/barsanuphe/radis)](http://goreportcard.com/report/barsanuphe/radis)
[![Coverage Status](https://coveralls.io/repos/github/barsanuphe/radis/badge.svg?branch=master)](https://coveralls.io/github/barsanuphe/radis?branch=master)
[![GPLv3](https://img.shields.io/badge/license-GPLv3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.en.html)

# Radis

## What it is

**radis** is a tool to help organize my music collection.
Can it be of any use to you?
Unless you like your music organized like I do, probably not.
See [Usage](#usage) to see what I mean.

So **radis** reads configuration files to move albums in the correct genre and
artist folders.
I could not get [beets](https://github.com/beetbox/beets) to organize things as
I wanted, so **radis** is what I use once [beets](https://github.com/beetbox/beets)
has imported (with correct tags, embedded art, etc) an album.

It does a few things:

- It helps with artists with more than one alias, so that an album is moved to
 `main alias/alias (year) album title`.
- It also helps sorting artists into the genres *I* have decided they belong to.
- If it is sorting them out of an "incoming" folder, it treats the albums as new
to your collection and automatically generates daily/monthly
[MPD](http://www.musicpd.org/) playlists of your new music.
- It can list and update the above-mentionned playlists if the albums have moved
since their creation.
- It can also tell you which albums are mp3s instead of flac.

And it's quite fast too.

*DISCLAIMER*: **radis** moves files around, deletes empty directories:
**You may lose data.**
Act as if you could lose everything and prepare accordingly.

## Why it is

[beets](https://github.com/beetbox/beets) can fetch genres from last.fm, but
most of the time the results feel wrong to me: too generic or too specific.

With **radis**, you decide which artist belong to which genre.
If you decide Beethoven isn't Death Metal after all, change the configuration
file and sync again.

Also, this is an excuse to learn Go.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Configuration examples](#configuration-example)

### Requirements

**Radis** runs on Linux (tested in Archlinux only).
It may run fine elsewhere, who knows.

It is written in Go, so everything fits into a nice executable.

### Installation

If you have a working Go environnment, go get it:

    $ go get -u github.com/barsanuphe/radis

The configuration files *radis_aliases.yaml* and *radis_genres.yaml* are
expected to be in `$XDG_CONFIG_HOME/radis/`.
You might want to `ln -s` your actual configuration files there.

### Usage

This command lists what was found in the configuration files:

    $ radis config show

This loads, sorts, and save the configuration files, for better readability:

    $ radis config save

This reorganizes your music collection in the `Root` indicated in `radis.yaml`:

    $ radis collection sync

Make sure `Root` is correct.
**radis** will stop if the path does not exist, but otherwise it will at least
delete empty directories in that `Root`.

Also, **radis** only works if you have this setup:

    root/
    |- Genre/
       |- Artist/
          |- Artist (year) Album

What if you have something different?
Then you should not use **radis**.

Before doing a `sync`, you can also try `radis collection check`, which will
just show what a `sync` would do.

Of course, you should only have flac versions of your music.
Sometimes they do not exist, so these albums have a `[MP3]` suffix in the folder
name.

To list those offending albums, and check you have not missed any:

    $ radis collection fsck

To list known playlists:

    $ radis playlist show

And to update them after a `radis sync` that has moved albums around:

    $ radis playlist update

When in doubt:

    $ radis help

### Configuration

**radis** uses three yaml files to describe
- useful paths
- artist aliases
- artist genre

`radis.yaml` looks like this:

    # where your music is, assumed to be your MPD music_directory
    Root:   /path/to/music/collection/
    # if albums are found there, they will be added to MPD playlists
    IncomingSubdir: INCOMING
    # albums without genres are put there
    UnsortedSubdir: UNCATEGORIZED
    # playlists are created there
    MPDPlaylistDirectory: /path/to/mpd/playlists/


`radis_aliases.yaml` looks like this:

    main_alias:
    - other alias
    - another one

`radis_genres.yaml` is not surprising either:

    genre:
    - artist
    - another one

Compilations should be in folders such as: `Various Artists (1937) Old Stuff`.
They can be associated to a genre in the yaml file:

    genre:
    - artist
    - Various Artists | compilation title

Remember you can use `radis config save` to reorder the files for aliases and
genres.

### Configuration examples

`radis_aliases.yaml`:

    MF DOOM:
    - MadVillain
    - JJ DOOM
    - Viktor Vaughn
    Radiohead:
    - Thom Yorke
    - Jonny Greenwood

`radis_genres.yaml` is not surprising either:

    Underground Hip-Hop:
    - MF DOOM
    Brit-Rock:
    - Radiohead
    - Blur
    Blues:
    - Various Artists | Rare Chicago Blues


