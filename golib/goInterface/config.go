// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package goInterface

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

	//	Books     []string
	//	Documents []string
	//	Downloads []string
	//	Movies    []string
	//	Music     []string
	//	Pictures  []string
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

/*
func (cfg *Config) AddBooks(dir string) {
	cfg.addDir("Books", dir)
}

func (cfg *Config) AddDocuments(dir string) {
	cfg.addDir("Documents", dir)
}

func (cfg *Config) AddDownloads(dir string) {
	cfg.addDir("Downloads", dir)
}

func (cfg *Config) AddMovies(dir string) {
	cfg.addDir("Movies", dir)
}

func (cfg *Config) AddMusic(dir string) {
	cfg.addDir("Music", dir)
}

func (cfg *Config) AddPictures(dir string) {
	cfg.addDir("Pictures", dir)
}
*/
