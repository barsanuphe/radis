# Radis

## What it is

**radis** is a tool to help organize my music collection.
It also is an excuse to learn Go.

It reads configuration files to move albums in the correct genre and artist
folders.
I could not get [beets](https://github.com/beetbox/beets) to organize things as
I wanted, so **radis** is what I use once [beets](https://github.com/beetbox/beets)
has imported (with correct tags, embedded art, etc) an album.

It does two things:

- It helps with artists with more than one alias, so that an album is moved to
 `main alias/alias (year) album title`.
- It also helps sorting artists into the genres *I* have decided they belong to.

[beets](https://github.com/beetbox/beets) can fetch genres from last.fm, but I
am not satisfied with the results, as the genres seem to be too generic or much
too specific to me.

It moves files around, deletes empty directories: **You may lose data.**
Act as if you could lose everything and prepare accordingly.

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

    $ go get XXXXXXX/radis
    $ go install ...radis

The configuration files *radis_aliases.yaml* and *radis_genres.yaml* are
expected to be in `$XDG_CONFIG_HOME/radis/`.
You might want to `ln -s` your actual configuration files there.

### Usage

This command lists what was found in the configuration files:

    $ radis show

This reorganizes your music collection in `FOLDER`:

    $ radis sync FOLDER

Make sure `FOLDER` is correct.
**radis** will stop if the path does not exist, but otherwise it will at least
delete empty directories in that `FOLDER`.

**radis** only works if you have this setup:

    root/
    |- Genre/
       |- Artist/
          |- Artist (year) Album

What if you have something different?
Then you should not use **radis**.

When in doubt:

    $ radis help

### Configuration

**radis** uses two yaml files to describe
- artist aliases
- artist genre

`radis_aliases.yaml` looks like this:

    main_alias:
    - other alias
    - another one

`radis_genres.yaml` is not surprising either:

    genre:
    - artist
    - another one

**radis** reorders the files each time it runs, so that they get easier to read.

### Configuration examples

`radis_aliases.yaml`:

    MF DOOM:
    - MadVillain
    - JJ DOOM
    Radiohead:
    - Thom Yorke
    - Jonny Greenwood

`radis_genres.yaml` is not surprising either:

    Underground Hip-Hop:
    - MF DOOM
    Brit-Rock:
    - Radiohead
    - Blur

