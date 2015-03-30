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

func makeConfig() *Config {
	cfg := &Config{}
	cfg.Devicename, _ = os.Hostname()
	cfg.AppFilesDir = "files"
	cfg.AddBooks(books)
	cfg.AddDocuments(documents)
	cfg.AddDownloads(downloads)
	cfg.AddMovies(movies)
	cfg.AddMusic(music)
	cfg.AddPictures(pictures)
	return cfg
}

func TestConfigAdds(t *testing.T) {
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
}

func TestSaveConfig(t *testing.T) {
	cfg := makeConfig()

	if err := cfg.Save(); err != nil {
		t.Error(err)
	}
}
