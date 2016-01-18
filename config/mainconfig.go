package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// TODO mark optional?

type MainConfig struct {
	Root                 string
	IncomingSubdir       string
	UnsortedSubdir       string
	MPDPlaylistDirectory string
}

func (mc *MainConfig) String() string {
	txt := "Radis configuration:\n"
	txt += "\tRoot: " + mc.Root + "\n"
	txt += "\tIncomingSubdir: " + mc.IncomingSubdir + "\n"
	txt += "\tUnsortedSubdir: " + mc.UnsortedSubdir + "\n"
	txt += "\tMPDPlaylistDirectory: " + mc.MPDPlaylistDirectory + "\n"
	return txt
}

func (mc *MainConfig) Check() (error) {
	// TODO: check all directories exist
	return nil
}

func (mc *MainConfig) Load(path string) (err error) {
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
		// TODO check that v exists, and that we have all keys!!!
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

func (mc *MainConfig) Write(path string) (err error) {
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
