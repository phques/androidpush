// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package goInterface

import (
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
)

func makeConfig(addDirs bool) *Config {
	cfg := &Config{}
	cfg.Devicename, _ = os.Hostname()
	cfg.AppFilesDir = "files"

	if addDirs {
		cfg.AddDir("Books", books)
		cfg.AddDir("Documents", documents)
		cfg.AddDir("Downloads", downloads)
		cfg.AddDir("Movies", movies)
		cfg.AddDir("Music", music)
		cfg.AddDir("Pictures", pictures)
	}
	return cfg
}

func TestConfigAddDir(t *testing.T) {
	cfg := makeConfig(false)
	cfg.AddDir("Books", "/home/philippe/Books")

	if len(cfg.Dirs) != 1 {
		t.Error("cfg.Dirs should be length 1")
	}

	dir, ok := cfg.Dirs["Books"]
	if !ok {
		t.Error("Books entry not found ok=false")
	} else if len(dir) != 1 || dir[0] != "/home/philippe/Books" {
		t.Error("Books entry list does not contain books dir")
	}
}

func _testConfigAdds(t *testing.T) {
	/*
		cfg := makeConfig()
		if cfg.Books[0] != books {
			t.Error("bad books")
		}
		if cfg.Documents[0] != documents {
			t.Error("bad documents")
		}
		if cfg.Downloads[0] != downloads {
			t.Error("bad downloads")
		}
		if cfg.Movies[0] != movies {
			t.Error("bad movies")
		}
		if cfg.Music[0] != music {
			t.Error("bad music")
		}
		if cfg.Pictures[0] != pictures {
			t.Error("bad pictures")
		}
	*/
}

func TestSaveConfig(t *testing.T) {
	cfg := makeConfig(true)

	if err := cfg.Save(); err != nil {
		t.Error(err)
	}
}
