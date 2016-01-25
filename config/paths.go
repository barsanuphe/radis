package config

import (
	"io/ioutil"

	"github.com/barsanuphe/radis/directory"
	"gopkg.in/yaml.v2"
)

// Paths contains useful paths for radis.
type Paths struct {
	Root                 string
	IncomingSubdir       string
	UnsortedSubdir       string
	MPDPlaylistDirectory string
}

func (mc *Paths) String() string {
	txt := "Radis configuration:\n"
	txt += "\tRoot: " + mc.Root + "\n"
	txt += "\tIncomingSubdir: " + mc.IncomingSubdir + "\n"
	txt += "\tUnsortedSubdir: " + mc.UnsortedSubdir + "\n"
	txt += "\tMPDPlaylistDirectory: " + mc.MPDPlaylistDirectory + "\n"
	return txt
}

// Check all the paths in MainConfig exist
func (mc *Paths) Check() (err error) {
	// check the required directories exist
	// the other directories can be created by radis
	if _, err := directory.GetExistingPath(mc.Root); err != nil {
		return err
	}
	if _, err := directory.GetExistingPath(mc.MPDPlaylistDirectory); err != nil {
		return err
	}
	return
}

// Load the configuration file where the paths are defined.
func (mc *Paths) Load(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	m := make(map[string]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	for k, v := range m {
		// TODO check that we have all keys!!!
		switch k {
		case "Root":
			mc.Root = v
		case "IncomingSubdir":
			mc.IncomingSubdir = v
		case "UnsortedSubdir":
			mc.UnsortedSubdir = v
		case "MPDPlaylistDirectory":
			mc.MPDPlaylistDirectory = v
		}
	}
	return
}

func (mc *Paths) Write(path string) (err error) {
	// TODO do better
	m := make(map[string]string)
	m["Root"] = mc.Root
	m["IncomingSubdir"] = mc.IncomingSubdir
	m["UnsortedSubdir"] = mc.UnsortedSubdir
	m["MPDPlaylistDirectory"] = mc.MPDPlaylistDirectory

	d, err := yaml.Marshal(&m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, d, 0777)
	return
}
