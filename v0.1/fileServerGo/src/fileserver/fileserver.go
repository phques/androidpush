// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// this will be the androidPush fileServer .. when I actually get to the real pgrm ! ;-p
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
	fmt.Printf("roots : %+v\n------\n", roots)
	fmt.Printf("roots : %#v\n------\n", roots)

	localPath := "/home/kwez/Pictures/07Aug/12Aout/pic1.jpg"

	foundRoot, found := roots.LookupLocal(localPath)
	if found {
		fmt.Println(localPath)
		fmt.Printf("Found: '%s' =>\n%s\n   %s\n",
			foundRoot.MediaRoot.Name, foundRoot.Base, foundRoot.Tail)
		fmt.Printf(" => %s/%s\n", foundRoot.MediaRoot.RemoteDirs[0], foundRoot.Tail)
	}
}
