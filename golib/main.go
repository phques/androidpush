// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"flag"
	"fmt"

	"github.com/phques/androidpush/golib/goInterface"
	_ "github.com/phques/androidpush/golib/goInterface/gen"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
	//	"net/http"
)

var standalone = flag.Bool("standalone", false, "running as a stand alone pgm")

func main() {
	fmt.Println("main")
	flag.Parse()
	app.Run(app.Callbacks{Start: start, Stop: stop})
}

func start() {
	fmt.Println("main.start")
	java.Init()

	if *standalone {
		goInterface.InitAppFilesDir("files")
		goInterface.Start()
	}
}

func stop() {
	fmt.Println("main.stop")
	if *standalone {
		//ps: not really required since the app is closing anyways ;-)
		goInterface.Stop()
	}
}
