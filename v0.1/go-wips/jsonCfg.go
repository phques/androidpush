package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Root directories for a media type, ie "music",
type MediaRoot struct {
	LocalDirs  []string // local root directories
	RemoteDirs []string // remote root directories
	Name       string   // media name ("music", "pictures"..)
	LocalIdx   int      // index of found local  root after lookup/selection
	RemoteIdx  int      // index of found remote root after lookup/selection
}

// index=name, ie roots["Music"] = MediaRoot{Name="Music"...}
type Roots map[string]MediaRoot

//-----------

// 'toString' for a MediaRoot
func (roots MediaRoot) String() string {
	return fmt.Sprintf("LocalDirs : %s, RemoteDirs: %s", roots.LocalDirs, roots.RemoteDirs)
}

//-----------

// Find a local root for 'dir' in root *MediaRoot
func (root *MediaRoot) LookupLocal(localDir string) (foundLocalIdx int, found bool) {
	// compare all lowercase (ie- Windows FS not case sensitive)
	lowerDir := strings.ToLower(localDir)
	var foundLen = 0 // use the longest root found

	for idx, localRoot := range root.LocalDirs {
		// if localRoot is a prefix of localDir then we found a root dir
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

//-----------

// Find a local root for 'dir' in Roots
func (roots Roots) LookupLocal(dir string) (foundRoot MediaRoot, found bool) {
	for _, mediaRoot := range roots {
		if localIdx, didFindLocal := mediaRoot.LookupLocal(dir); didFindLocal {
			// found a root dir, copy & save
			foundRoot = mediaRoot
			foundRoot.LocalIdx = localIdx
			found = true
			break
		}
	}
	return
}

//-----------

// Create a Roots map from a json text file config (list of MediaRoot)
func readRootdirs(jsonCfgFile string) (roots Roots, err error) {
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

//-----------

func main() {
	// read media roots cfg
	roots, err := readRootdirs("config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(roots)

	// test Lookup a local dir
	var foundRoot MediaRoot
	var found bool
	foundRoot, found = roots.LookupLocal("/home/kwez/Pictures/12Aout")
	if found {
		fmt.Printf("Found: '%s' => %s @ %d\n",
			foundRoot.Name, foundRoot, foundRoot.LocalIdx)
	}

}
