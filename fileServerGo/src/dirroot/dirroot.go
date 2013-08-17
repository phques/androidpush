// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package dirroot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

//---------

// Root directories for a media type, ie "music",
type MediaRoot struct {
	LocalDirs  []string // local root directories
	RemoteDirs []string // remote root directories
	Name       string   // media name ("music", "pictures"..)
	LocalIdx   int
}

// index=name, ie roots["Music"] = MediaRoot{Name="Music"...}
type Roots map[string]MediaRoot

// Result of a local root lookup
type FoundLocalRoot struct {
	FullPath  string
	MediaRoot MediaRoot
	Idx       int // index of found local root
	Base      string
	Tail      string
}

//-----------

// 'toString' for a MediaRoot
func (roots MediaRoot) String() string {
	return fmt.Sprintf("LocalDirs: %s, RemoteDirs: %s", roots.LocalDirs, roots.RemoteDirs)
}

//-----------

// Find a local root for 'dir' in root *MediaRoot
func (root *MediaRoot) _LookupLocal(localPath string) (foundLocalIdx int, found bool) {
	// compare all lowercase (ie- Windows FS not case sensitive)
	lowerDir := strings.ToLower(localPath)
	var foundLen = 0 // use the longest root found

	for idx, localRoot := range root.LocalDirs {
		// if localRoot is a prefix of localPath then we found a root dir
		if strings.HasPrefix(lowerDir, strings.ToLower(localRoot)) {
			// only use this one if it is a longer prefix than prev found
			if len(localRoot) > foundLen {
				foundLocalIdx = idx
				found = true
				foundLen = len(localRoot)
			}
		}
	}
	return
}

// Find a local root for 'dir' in root *MediaRoot
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

// Find a local root for 'dir' in Roots
func (roots Roots) _LookupLocal(dir string) (foundRoot MediaRoot, found bool) {
	for _, mediaRoot := range roots {
		if localIdx, didFindLocal := mediaRoot._LookupLocal(dir); didFindLocal {
			// found a root dir, copy & save
			foundRoot = mediaRoot
			foundRoot.LocalIdx = localIdx
			found = true
			break
		}
	}
	return
}

// Find a local root for 'dir' in Roots
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

// Create a Roots map from a json text file config (list of MediaRoot)
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
