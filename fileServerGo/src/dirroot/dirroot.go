// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// Package dirroot provides types & methods to identify 'root' directories for filepaths.
// The roots (local & corresponding remote) are read from a JSON config file
package dirroot

import (
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"strings"
)

//---------

// Local (& remote) 'root' directories for a media type, ie "music",
type MediaRoot struct {
	LocalDirs  []string // local root directories
	RemoteDirs []string // remote root directories
	Name       string   // media name ("music", "pictures"..)
}

// Map of MediaRoot indexed by media name
//  ie roots["Music"] = MediaRoot{Name="Music"...}
type Roots map[string]MediaRoot

// FoundLocalRoot is the result of a local root lookup
type FoundLocalRoot struct {
	MediaRoot MediaRoot
	Idx       int // index of found local root MediaRoot.LocalDirs[Idx]
	FullPath  string
	Base      string
	Tail      string
}

//-----------

// LookupLocal finds a local root for 'dir' in root *MediaRoot
func (root *MediaRoot) LookupLocal(localPath string) (foundRoot FoundLocalRoot, found bool) {
	// compare all lowercase (ie- Windows FS not case sensitive)
	lowerDir := strings.ToLower(localPath)
	var foundLen = 0 // use the longest root found

	for idx, localRoot := range root.LocalDirs {

		// if localRoot is a prefix of localPath then we found a local root dir
		if strings.HasPrefix(lowerDir, strings.ToLower(localRoot)) {

			// only use this one if it is a longer prefix than prev found
			if len(localRoot) > foundLen {
				foundLen = len(localRoot)

				foundRoot = FoundLocalRoot{
					FullPath:  localPath,
					MediaRoot: *root,
					Idx:       idx,
					Base:      localRoot,
					Tail:      localPath[len(localRoot)+1:], // skip leading '/'
				}
				found = true
			}
		}
	}
	return
}

//-----------

// LookupLocal finds a local root for 'dir' in Roots
func (roots Roots) LookupLocal(dir string) (foundRoot FoundLocalRoot, found bool) {
	for _, mediaRoot := range roots {
		// try to a root dir of 'dir' in one of mediaRoot.Locals
		foundRoot, found = mediaRoot.LookupLocal(dir)
		if found { // found a root dir
			break
		}
	}
	return
}

//-----------

// ReadRootdirs creates a Roots map from a json text file config (list of MediaRoot)
func ReadRootdirs(jsonCfgFile string) (roots Roots, err error) {
	// read json config file as text
	jsonCfg, err := ioutil.ReadFile(jsonCfgFile)
	if err != nil {
		return
	}

	// unmarshall json to []MediaRoot (is a json list)
	var mediaRootsArr []MediaRoot
	err = json.Unmarshal([]byte(jsonCfg), &mediaRootsArr)
	if err != nil {
		return
	}

	// convert []MediaRoot to Roots
	roots = Roots{}
	for _, mediaRoot := range mediaRootsArr {
		roots[mediaRoot.Name] = mediaRoot
	}

	return
}
