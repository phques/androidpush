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
	Books     []string
	Documents []string
	Downloads []string
	Movies    []string
	Music     []string
	Pictures  []string
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

func (cfg *Config) AddBooks(dir string) {
	cfg.Books = append(cfg.Books, dir)
}

func (cfg *Config) AddDocuments(dir string) {
	cfg.Documents = append(cfg.Documents, dir)
}

func (cfg *Config) AddDownloads(dir string) {
	cfg.Downloads = append(cfg.Downloads, dir)
}

func (cfg *Config) AddMovies(dir string) {
	cfg.Movies = append(cfg.Movies, dir)
}

func (cfg *Config) AddMusic(dir string) {
	cfg.Music = append(cfg.Music, dir)
}

func (cfg *Config) AddPictures(dir string) {
	cfg.Pictures = append(cfg.Pictures, dir)
}
