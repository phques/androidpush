// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package gopush

import (
	"encoding/json"
	"io/ioutil"
)

// Config contains the configuration info for androidPush
type Config struct {
	Devicename  string // reported name to mppq service query responses
	AppFilesDir string // app's files dir, we store config file there

	// 'root' directories
	// ie ["Books"]("~/Books", "~/Calibre/lib/books")
	Dirs map[string][]string
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

// Save save the config to file in JSON format
func (cfg *Config) Save(filename string) error {
	//  serialize config to JSON
	jsonStr, err := json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return err
	}

	// save to file
	err = ioutil.WriteFile(filename, []byte(jsonStr), 0664)
	if err != nil {
		return err
	}

	return nil
}

// Load reads a Config from saved file, JSON format
func (cfg *Config) Load(filename string) error {
	// read file
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// deserialize data into cfg
	err = json.Unmarshal(jsonStr, cfg)
	if err != nil {
		return err
	}

	return nil
}
