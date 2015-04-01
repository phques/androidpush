// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package gopush

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	books     string = "/home/philippe/Books"
	documents string = "/home/philippe/Documents"
	downloads string = "/home/philippe/Downloads"
	movies    string = "/home/philippe/Movies"
	music     string = "/home/philippe/Music"
	pictures  string = "/home/philippe/Pictures"

	testFilename = "config_test.json"
)

// MakeConfig creates & fills a Config
func MakeConfig() *Config {
	cfg := &Config{}
	cfg.Devicename, _ = os.Hostname()
	cfg.AppFilesDir = "files"

	cfg.AddDir("Books", books)
	cfg.AddDir("Documents", documents)
	cfg.AddDir("Downloads", downloads)
	cfg.AddDir("Movies", movies)
	cfg.AddDir("Music", music)
	cfg.AddDir("Pictures", pictures)

	return cfg
}

// helper for TestConfigAddDir
func (cfg *Config) verifyAddedDir(dirName, dirValue string, nbEntries int, t *testing.T) {
	// check entry was added
	if len(cfg.Dirs) != nbEntries {
		t.Error("cfg.Dirs should be length 1")
	}

	// fetch & check added entry
	dir, ok := cfg.Dirs[dirName]
	if !ok {
		t.Error(dirName + " entry not found, ok=false")
	} else if len(dir) != 1 || dir[0] != dirValue {
		t.Error(dirName + " entry list does not equal value " + dirValue)
	}

}

// Test adding dir entries to a Config
func TestConfigAddDir(t *testing.T) {
	// create empty Config, add & test dir entries
	cfg := &Config{}

	cfg.AddDir("Books", books)
	cfg.verifyAddedDir("Books", books, 1, t)

	cfg.AddDir("Music", music)
	cfg.verifyAddedDir("Music", music, 2, t)
}

func TestSaveLoadConfig(t *testing.T) {
	// create a Config & save to file
	cfg := MakeConfig()
	err := cfg.Save(testFilename)
	if err != nil {
		t.Error(err)
	}

	// Create empty Config & load from file
	cfgBack := &Config{}
	err = cfgBack.Load(testFilename)
	if err != nil {
		t.Error(err)
	}

	// Convert both to json to compare equality
	cfgStr, err1 := json.Marshal(*cfg)
	cfgBackStr, err2 := json.Marshal(*cfgBack)

	if err1 != nil || err2 != nil {
		t.Error("failed to marshal ", err1, err2)
	}

	if string(cfgBackStr) != string(cfgStr) {
		t.Error("saved & loaded Configs are not equal")
	}
}
