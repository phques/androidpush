// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"dirroot"
	"fmt"
)

//-----------

func main() {
	// read media roots cfg
	roots, err := dirroot.ReadRootdirs("config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(roots)

	// test Lookup a local dir
	//var foundRoot dirroot.MediaRoot
	//var found bool
	//foundRoot, found = roots.LookupLocal("/home/kwez/Pictures/07Aug/12Aout/pic1.jpg")
	//if found {
	//	fmt.Printf("Found: '%s' => %s @ %d\n",
	//		foundRoot.Name, foundRoot, foundRoot.LocalIdx)
	//}
	var foundRoot dirroot.FoundLocalRoot
	var found bool
	foundRoot, found = roots.LookupLocal("/home/kwez/Pictures/07Aug/12Aout/pic1.jpg")
	if found {
		fmt.Printf("Found: '%s' =>\n%s\n   %s\n",
			foundRoot.MediaRoot.Name, foundRoot.Base, foundRoot.Tail)
	}

}
