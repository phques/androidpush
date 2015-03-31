// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package gopush

import (
	"encoding/json"
	"log"
)

type Config struct {
	Devicename  string // reported name to mppq service query responses
	AppFilesDir string // app's files dir, we store config file there

	// config file directories, used to populate config file 1st time
	// ie ["Books"]("~/Books", "~/Calibre/lib/books")
	Dirs map[string][]string
}

func (cfg *Config) Save() error {

	json, err := json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return err
	}
	//##debug
	log.Printf("%v", string(json))
	return nil
}

// AddDir adds entry dir in cfg.Dirs[name]
func (cfg *Config) AddDir(name string, dir string) {
	// create Dirs if nil
	if cfg.Dirs == nil {
		cfg.Dirs = make(map[string][]string)
	}
	// Get current list of directories for [name] and add new dir
	dirs := cfg.Dirs[name]
	cfg.Dirs[name] = append(dirs, dir)

}
