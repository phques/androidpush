// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// package dirroot tests
package dirroot

import (
	//"fmt"
	"testing"
)

var roots = Roots{
	"Music":     MediaRoot{LocalDirs: []string{"/home/kwez/Music"}, RemoteDirs: []string{"/storage/extSdCard/Music"}, Name: "Music"},
	"Downloads": MediaRoot{LocalDirs: []string{"/home/kwez/Downloads"}, RemoteDirs: []string{"/storage/extSdCard/Download"}, Name: "Downloads"},
	"Pictures":  MediaRoot{LocalDirs: []string{"/home/kwez/Pictures"}, RemoteDirs: []string{"/storage/extSdCard/Pictures"}, Name: "Pictures"},
	"Movies":    MediaRoot{LocalDirs: []string{"/home/kwez/Videos"}, RemoteDirs: []string{"/storage/extSdCard/Movies"}, Name: "Movies"},
}

func TestLookupLocal(t *testing.T) {
	// try to find a local dir root for this local path
	localPath := "/home/kwez/Pictures/07Aug/12Aout/pic1.jpg"

	// call LookupLocal
	foundRoot, found := roots.LookupLocal(localPath)
	_ = foundRoot
	if !found {
		t.Errorf("failed to find local root dir '%s'", localPath)
	}
}
